package watchlog

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Watcher struct {
	sync.Mutex

	client apiclient.ChainHttpClient

	db            store.DB
	factoryBatch  store.Batch
	transferBatch store.Batch
	commitDelay   time.Duration // delay to commit batch

	watchingNfts       []common.Address
	originWatchNftsCtx context.Context
	cancelWatchNfts    context.CancelFunc

	errCh chan error
}

// NOTE: delay is millisecond
func NewWatcher(db store.DB, client apiclient.ChainHttpClient, delay int) *Watcher {
	return &Watcher{
		client:        client,
		db:            db,
		factoryBatch:  db.Batch(),
		transferBatch: db.Batch(),
		commitDelay:   time.Duration(delay) * time.Millisecond,
		errCh:         make(chan error),
	}
}

func (w *Watcher) Close() (err error) {
	w.Lock()
	defer w.Unlock()

	defer close(w.errCh)
	defer w.factoryBatch.Close()
	defer w.transferBatch.Close()

	if !w.factoryBatch.Empty() {
		err = w.factoryBatch.Commit()
		// not return here, because we need to commit transfer batch
	}

	if !w.transferBatch.Empty() {
		err = w.transferBatch.Commit()
	}

	return
}

func (w *Watcher) Err() <-chan error {
	return w.errCh
}

func (w *Watcher) Start(ctx context.Context, factoryAddress common.Address) (err error) {
	if err = w.startWatchingTransferLog(ctx); err != nil {
		return fmt.Errorf("failed to start watching transfer log: %w", err)
	}
	if err = w.startWatchingFactoryLog(ctx, factoryAddress); err != nil {
		return fmt.Errorf("failed to start watching factory log: %w", err)
	}

	return
}

func (w *Watcher) startWatchingFactoryLog(ctx context.Context, address common.Address) error {
	go func() {
		if err := w.client.WatchFactoryLog(ctx, address, w.handleFactoryEvent); err != nil {
			w.errCh <- fmt.Errorf("failed to watch factory log: %w", err)
		}
	}()

	return nil
}

func (w *Watcher) startWatchingTransferLog(ctx context.Context) error {
	nfts, err := w.allNFTs()
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to get all nfts: %w", err)
	}

	// the total number of watching nfts will increase, thus we need to allocate more space
	margin := 64
	w.watchingNfts = make([]common.Address, len(nfts)+margin)
	for i := range nfts {
		w.watchingNfts[i] = common.HexToAddress(nfts[i].Address)
	}

	// wrap context to cancel watching nfts
	w.originWatchNftsCtx = ctx

	go func() {
		ctx, w.cancelWatchNfts = context.WithCancel(w.originWatchNftsCtx)

		w.Lock()
		targets := w.watchingNfts
		w.Unlock()

		if err := w.client.WatchTransferLog(ctx, targets, w.handleTransferEvent); err != nil {
			w.errCh <- fmt.Errorf("failed to watch transfer log: %w", err)
		}
	}()

	return nil
}

func (w *Watcher) handleFactoryEvent(e *factory.FactoryNFTCreated) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		nft = data.NFTContract{Address: e.Nft.Hex()}
		now = timestamppb.Now()
		wg  = sync.WaitGroup{}
	)
	if err = w.client.FetchNFTInfo(ctx, &nft); err != nil {
		return fmt.Errorf("failed to fetch nft info: %w", err)
	}

	// store and new nft contract needs to be watch
	w.Lock()
	w.factoryBatch.PutWithTime(now, &nft)
	w.watchingNfts = append(w.watchingNfts, e.Nft)
	w.Unlock()

	defer func() {
		w.Lock()
		len := w.factoryBatch.Len()
		w.Unlock()

		// commit the batch `commitDelay` after the first event is received
		if len == 1 {
			time.Sleep(w.commitDelay)

			w.Lock()
			defer w.Unlock()

			if err = w.factoryBatch.Commit(); err != nil {
				err = fmt.Errorf("failed to commit factory batch: %w", err)
				w.errCh <- err
			}
		}
	}()

	wg.Add(1)
	oldCancel := w.cancelWatchNfts

	// restart watching transfer log with new nft contract
	go func() {
		wg.Done()
		ctx, w.cancelWatchNfts = context.WithCancel(w.originWatchNftsCtx)

		w.Lock()
		targets := w.watchingNfts
		w.Unlock()

		if err = w.client.WatchTransferLog(ctx, targets, w.handleTransferEvent); err != nil {
			w.errCh <- fmt.Errorf("failed to start watching new transfer log: %w", err)
		}
	}()

	// wait to cancel old nfts watching until the new watching nft contract is started
	wg.Wait()
	oldCancel()

	return
}

func (w *Watcher) handleTransferEvent(e *ierc721.ContractTransfer) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		token *data.Token
		now   = timestamppb.Now()
	)

	// prepare token
	if token, err = w.getToken(e.Raw.Address.Hex(), e.TokenId.String(), now); err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	token.Owner = e.To.Hex()
	if token.Meta == nil {
		// initialize meta
		if err = w.client.GetTokenMeta(ctx, token); err != nil {
			return fmt.Errorf("failed to fetch token info: %w", err)
		}
	}

	// prepare history and owner index
	var (
		history = &data.TransferHistory{
			Address:         e.Raw.Address.Hex(),
			TokenId:         e.TokenId.String(),
			From:            e.From.Hex(),
			To:              e.To.Hex(),
			BlockNumber:     e.Raw.BlockNumber,
			IndexLogInBlock: uint32(e.Raw.Index),
		}
		ownerIndex = data.NewTokenOwnerIndex(token)
	)

	// store
	w.Lock()
	w.transferBatch.PutWithTime(now, token, history, ownerIndex)
	w.Unlock()

	defer func() {
		w.Lock()
		len := w.factoryBatch.Len()
		w.Unlock()

		// commit the batch `commitDelay` after the first event is received
		if len == 1 {
			time.Sleep(w.commitDelay)

			w.Lock()
			defer w.Unlock()

			if err = w.transferBatch.Commit(); err != nil {
				err = fmt.Errorf("failed to commit transfer batch: %w", err)
				w.errCh <- err
			}
		}
	}()

	return
}

func (w *Watcher) allNFTs() ([]*data.NFTContract, error) {
	results := []*data.NFTContract{{}}
	if err := w.db.List(data.PrefixNFTContract, &results); err != nil {
		return nil, fmt.Errorf("failed to list nft contracts: %w", err)
	}
	return results, nil
}

func (w *Watcher) getToken(address, tokenId string, now *timestamppb.Timestamp) (*data.Token, error) {
	t := data.Token{Address: address, TokenId: tokenId}
	if err := w.db.Get(t.Key(), &t); err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return &t, fmt.Errorf("failed to get last factory log fetched: %w", err)
		} else {
			// first time to fetch
			t.CreatedAt = now
		}
	}
	return &t, nil
}
