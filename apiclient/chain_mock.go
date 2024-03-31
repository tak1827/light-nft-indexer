package apiclient

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
)

var _ ChainHttpClient = (*MockChainClient)(nil)

// TODO: add delay of calling rpc
type MockChainClient struct {
	sync.Mutex

	fe      []*factory.FactoryNFTCreated
	feIndex int
	te      []*ierc721.ContractTransfer
	teIndex int

	endHeight *uint64

	delay time.Duration
	errCh chan error
}

func NewMockChainClient(fe []*factory.FactoryNFTCreated, te []*ierc721.ContractTransfer, endHeight *uint64, delay int) (c MockChainClient) {
	c.fe = fe
	c.te = te
	c.endHeight = endHeight
	c.delay = time.Duration(delay) * time.Millisecond
	c.errCh = make(chan error)
	return
}

func (c *MockChainClient) EmitErr(err error) {
	c.errCh <- err
}

func (c *MockChainClient) FetchFactoryLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*factory.FactoryNFTCreated, nextStart uint64, err error) {
	return c.fe, *c.endHeight, nil
}

func (c *MockChainClient) FetchTransferLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*ierc721.ContractTransfer, nextStart uint64, err error) {
	return c.te, *c.endHeight, nil
}

func (c *MockChainClient) FetchNFTInfo(ctx context.Context, d *data.NFTContract) (err error) {
	d.Name = "name"
	d.Symbol = "symbol"
	d.TotalSupply = 100
	return
}

func (c *MockChainClient) GetTokenMeta(ctx context.Context, d *data.Token) (err error) {
	d.Meta = &data.TokenMeta{
		Origin: "meta",
		Image: &data.TokenMetaImage{
			Type: data.ImageType_IMAGE_TYPE_UNSPECIFIED,
			Data: "data",
		},
	}
	return
}

func (c *MockChainClient) WatchFactoryLog(ctx context.Context, address common.Address, callback func(*factory.FactoryNFTCreated) error) error {
	ticker := time.NewTicker(c.delay)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-c.errCh:
			return err
		case <-ticker.C:
			if len(c.fe) <= c.feIndex {
				ticker.Stop()
				break
			}
			if err := callback(c.fe[c.feIndex]); err != nil {
				return err
			}
			c.feIndex++
		}
	}
}

func (c *MockChainClient) WatchTransferLog(ctx context.Context, addresses []common.Address, callback func(*ierc721.ContractTransfer) error) error {

	ticker := time.NewTicker(c.delay)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-c.errCh:
			return err
		case <-ticker.C:
			c.Lock()
			i := c.teIndex
			c.teIndex++
			c.Unlock()

			if len(c.te) <= i {
				ticker.Stop()
				break
			}

			// ensure the address of event is in the list
			for _, address := range addresses {
				if address.Cmp(c.te[i].Raw.Address) != 0 {
					break
				}
			}

			if err := callback(c.te[i]); err != nil {
				return err
			}
		}
	}
}
