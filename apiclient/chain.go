package apiclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/tak1827/go-cache/lru"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
)

const (
	DefaultCacheSize = 64
	DefaultCacheTTL  = 60 * 60 * 24 // 1day
	MaxRetrivalLogs  = 256
)

type ChainHttpClient struct {
	c      *ethclient.Client
	logger zerolog.Logger

	// cache contract instance
	cache lru.LRUCache
	// reuse filter options
	filterOpts bind.FilterOpts
}

func NewChainClient(ctx context.Context, endpoint string, logger zerolog.Logger) (c ChainHttpClient, err error) {
	rpcclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		err = fmt.Errorf("failed to create rpc client: %w", err)
		return
	}
	c.c = ethclient.NewClient(rpcclient)

	// try to connect
	if _, err = c.c.NetworkID(ctx); err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			err = fmt.Errorf("%w: %w", ErrConnectRefused, err)
		}
		err = fmt.Errorf("failed to connect to endpoint(%s): %w", endpoint, err)
		return
	}

	c.logger = logger
	c.cache = lru.NewCache(DefaultCacheSize, DefaultCacheTTL)
	c.filterOpts = bind.FilterOpts{}

	return
}

func (c *ChainHttpClient) FetchTransferLog(ctx context.Context, address common.Address, startHeight, endHeight uint64) (events []*ierc721.ContractTransfer, nextStart uint64, err error) {
	var (
		contract       *ierc721.Contract
		isExpectedType bool
	)
	value, exist := c.cache.Get(address.Hex())

	if exist {
		contract, isExpectedType = value.(*ierc721.Contract)
	}

	if !exist || isExpectedType {
		if contract, err = ierc721.NewContract(address, c.c); err != nil {
			err = fmt.Errorf("failed to create contract: %w", err)
		}
		c.cache.Add(address.Hex(), contract)
	}

	c.setFilterOpts(ctx, startHeight, endHeight)
	itr, err := contract.FilterTransfer(&c.filterOpts, nil, nil, nil)
	defer itr.Close()

	for itr.Next() {
		if MaxRetrivalLogs < len(events) {
			nextStart = itr.Event.Raw.BlockNumber
			break
		}
		events = append(events, itr.Event)
	}

	if err = itr.Error(); err != nil {
		err = fmt.Errorf("failed to iterate logs: %w", err)
		return
	}

	if nextStart == 0 {
		nextStart = endHeight
	}
	return
}

func (c *ChainHttpClient) setFilterOpts(ctx context.Context, startHeight, endHeight uint64) {
	c.filterOpts.Start = startHeight
	c.filterOpts.End = &endHeight
	c.filterOpts.Context = ctx
}
