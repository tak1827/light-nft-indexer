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
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
)

const (
	DefaultCacheSize = 64
	DefaultCacheTTL  = 60 * 60 * 24 // 1day
	MaxRetrivalLogs  = 256
)

var _ ChainHttpClient = (*EthHttpClient)(nil)

type EthHttpClient struct {
	c      *ethclient.Client
	logger zerolog.Logger
	cache  lru.LRUCache // cache contract instance
}

func NewChainClient(ctx context.Context, endpoint string, logger zerolog.Logger) (c EthHttpClient, err error) {
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

	return
}

func (c *EthHttpClient) FetchTransferLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*ierc721.ContractTransfer, nextStart uint64, err error) {
	contract, err := c.getNFTContract(ctx, address)
	if err != nil {
		err = fmt.Errorf("failed to get contract: %w", err)
		return
	}

	if endHeight == nil {
		*endHeight, err = c.c.BlockNumber(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get latest block number: %w", err)
			return
		}
	}

	opts := bind.FilterOpts{
		Start:   startHeight,
		End:     endHeight,
		Context: ctx,
	}
	itr, err := contract.FilterTransfer(&opts, nil, nil, nil)
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
		nextStart = *endHeight
	}
	return
}

func (c *EthHttpClient) FetchFactoryLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*factory.FactoryNFTCreated, nextStart uint64, err error) {
	var (
		contract       *factory.Factory
		isExpectedType bool
	)
	value, exist := c.cache.Get(address.Hex())

	if exist {
		contract, isExpectedType = value.(*factory.Factory)
	}

	if !exist || isExpectedType {
		if contract, err = factory.NewFactory(address, c.c); err != nil {
			err = fmt.Errorf("failed to create contract: %w", err)
		}
		c.cache.Add(address.Hex(), contract)
	}

	if endHeight == nil {
		*endHeight, err = c.c.BlockNumber(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get latest block number: %w", err)
			return
		}
	}

	opts := bind.FilterOpts{
		Start:   startHeight,
		End:     endHeight,
		Context: ctx,
	}
	itr, err := contract.FilterNFTCreated(&opts, nil)
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
		nextStart = *endHeight
	}
	return
}

func (c *EthHttpClient) FetchNFTInfo(ctx context.Context, d *data.NFTContract) (err error) {
	if d.Address == "" {
		err = fmt.Errorf("address is empty")
		return
	}

	contract, err := c.getNFTContract(ctx, common.HexToAddress(d.Address))
	if err != nil {
		err = fmt.Errorf("failed to get contract: %w", err)
		return
	}

	opt := bind.CallOpts{Context: ctx}
	if d.Name, err = contract.Name(&opt); err != nil {
		err = fmt.Errorf("failed to get name: %w", err)
		return
	}
	if d.Symbol, err = contract.Symbol(&opt); err != nil {
		err = fmt.Errorf("failed to get symbol: %w", err)
		return
	}
	if supply, err := contract.TotalSupply(&opt); err == nil {
		// ignore error, only update if succeeded
		d.TotalSupply = supply.Uint64()
	}
	return
}

func (c *EthHttpClient) getNFTContract(ctx context.Context, address common.Address) (contract *ierc721.Contract, err error) {
	value, exist := c.cache.Get(address.Hex())

	if exist {
		contract, _ = value.(*ierc721.Contract)
		return
	}

	if contract, err = ierc721.NewContract(address, c.c); err != nil {
		err = fmt.Errorf("failed to create contract: %w", err)
		return
	}
	c.cache.Add(address.Hex(), contract)
	return
}
