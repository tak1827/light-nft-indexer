package fetchlog

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Fetcher struct {
	db     store.DB
	client apiclient.ChainHttpClient
	batch  store.Batch
}

func NewFetcher(db store.DB, client apiclient.ChainHttpClient, batch store.Batch) *Fetcher {
	return &Fetcher{
		db:     db,
		client: client,
		batch:  batch,
	}
}

func (f *Fetcher) FetchAll(ctx context.Context, address common.Address, now *timestamppb.Timestamp) error {
	nfts, err := f.FetchNFTs(ctx, address, now, true)
	if err != nil {
		return fmt.Errorf("failed to fetch nft: %w", err)
	}

	for i := range nfts {
		if _, err = f.FetchTokens(ctx, address, now, &nfts[i], true); err != nil {
			return fmt.Errorf("failed to fetch token: %w", err)
		}
	}

	return nil
}

func (f *Fetcher) FetchNFTs(ctx context.Context, address common.Address, now *timestamppb.Timestamp, withCommit bool) ([]data.NFTContract, error) {
	b, err := f.getBlock(data.BlockType_BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED, "", now)
	if err != nil {
		return nil, err
	}

	events, nextStart, err := f.client.FetchFactoryLog(ctx, address, b.Height, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch factory log: %w", err)
	}

	nfts := make([]data.NFTContract, len(events))
	for i := range events {
		nfts[i] = data.NFTContract{Address: events[i].Nft.Hex()}
		if err = f.client.FetchNFTInfo(ctx, &nfts[i]); err != nil {
			return nil, fmt.Errorf("failed to fetch nft info: %w", err)
		}
		// store the fetched events
		f.batch.PutWithTime(now, &nfts[i])
	}

	if withCommit {
		if err = f.commit(b, nextStart, now); err != nil {
			return nil, err
		}
	}

	return nfts, nil
}

func (f *Fetcher) FetchTokens(ctx context.Context, address common.Address, now *timestamppb.Timestamp, nft *data.NFTContract, withCommit bool) ([]*data.Token, error) {
	b, err := f.getBlock(data.BlockType_BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED, nft.Address, now)
	if err != nil {
		return nil, err
	}

	events, nextStart, err := f.client.FetchTransferLog(ctx, address, b.Height, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transfer log: %w", err)
	}

	tokens := make([]*data.Token, len(events))
	for i := range events {
		if tokens[i], err = f.getToken(nft.Address, events[i].TokenId.String(), now); err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
		tokens[i].Owner = events[i].To.Hex()
		tokens[i].TransferHistories = append(tokens[i].TransferHistories, data.NewTransferHistory(nft.Address, events[i].TokenId.String(), events[i].From.Hex(), events[i].To.Hex(), now))
		if tokens[i].Meta == nil {
			// initialize meta
			if err = f.client.GetTokenMeta(ctx, tokens[i]); err != nil {
				return nil, fmt.Errorf("failed to fetch token info: %w", err)
			}
		}
		f.batch.PutWithTime(now, tokens[i])
	}

	if withCommit {
		if err = f.commit(b, nextStart, now); err != nil {
			return nil, err
		}
	}

	return tokens, nil
}

func (f *Fetcher) getBlock(t data.BlockType, sub string, now *timestamppb.Timestamp) (*data.Block, error) {
	b := data.Block{Type: t, SubIdentifier: sub}
	if err := f.db.Get(b.Key(), &b); !errors.Is(err, store.ErrNotFound) {
		return &b, fmt.Errorf("failed to get last factory log fetched: %w", err)
	} else {
		// first time to fetch
		b.Height = 0
		b.CreatedAt = now
	}

	return &b, nil
}

func (f *Fetcher) getToken(address, tokenId string, now *timestamppb.Timestamp) (*data.Token, error) {
	t := data.Token{Address: address, TokenId: tokenId}
	if err := f.db.Get(t.Key(), &t); !errors.Is(err, store.ErrNotFound) {
		return &t, fmt.Errorf("failed to get last factory log fetched: %w", err)
	} else {
		// first time to fetch
		t.CreatedAt = now
	}

	return &t, nil
}

func (f *Fetcher) commit(b *data.Block, nextStart uint64, now *timestamppb.Timestamp) error {
	// update the last fetched height
	b.Height = nextStart
	f.batch.PutWithTime(now, b)

	// commit batch for a moment
	if err := f.batch.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch: %w", err)
	}

	// reset batch, to reuse
	f.batch.Reset()

	return nil
}
