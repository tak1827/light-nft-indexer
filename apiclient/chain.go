package apiclient

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// TODO: add delay of calling rpc
type EthHttpClient struct {
	c      *ethclient.Client
	logger zerolog.Logger
	cache  lru.LRUCache // cache contract instance

	nftAbi *abi.ABI
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

	if c.nftAbi, err = factory.FactoryMetaData.GetAbi(); err != nil {
		err = fmt.Errorf("failed to parse abi: %w", err)
		return
	}

	return
}

func (c *EthHttpClient) FetchFactoryLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*factory.FactoryNFTCreated, nextStart uint64, err error) {
	contract, err := c.getFactoryContract(ctx, address)
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

func (c *EthHttpClient) GetTokenMeta(ctx context.Context, d *data.Token) (err error) {
	if d.Address == "" || d.TokenId == "" {
		err = fmt.Errorf("address or token id is empty")
		return
	}

	contract, err := c.getNFTContract(ctx, common.HexToAddress(d.Address))
	if err != nil {
		err = fmt.Errorf("failed to get contract: %w", err)
		return
	}

	// initialize meta
	if d.Meta == nil {
		d.Meta = &data.TokenMeta{}
	}

	var (
		opt     = bind.CallOpts{Context: ctx}
		tokenId = new(big.Int)
		success bool
	)
	if tokenId, success = tokenId.SetString(d.TokenId, 10); !success {
		err = fmt.Errorf("failed to convert token id(%s) to big.Int", d.TokenId)
		return
	}
	if d.Meta.Origin, err = contract.TokenURI(&opt, tokenId); err != nil {
		err = fmt.Errorf("failed to get name: %w", err)
		return
	}
	return
}

func (c *EthHttpClient) WatchFactoryLog(ctx context.Context, address common.Address, callback func(*factory.FactoryNFTCreated) error) error {
	contract, err := c.getFactoryContract(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get contract: %w", err)
	}

	var (
		ch   = make(chan *factory.FactoryNFTCreated, 128)
		opts = bind.WatchOpts{
			Context: ctx,
		}
	)

	sub, err := contract.WatchNFTCreated(&opts, ch, nil)
	if err != nil {
		return fmt.Errorf("failed to start subscription: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err = <-sub.Err():
			if err != nil {
				return fmt.Errorf("failed to watch transfer: %w", err)
			}
		case event := <-ch:
			if err = callback(event); err != nil {
				return fmt.Errorf("failed to callback: %w", err)
			}
		}
	}
}

func (c *EthHttpClient) WatchTransferLog(ctx context.Context, addresses []common.Address, callback func(*ierc721.ContractTransfer) error) error {
	topics, err := abi.MakeTopics([]interface{}{c.nftAbi.Events["Transfer"].ID})
	if err != nil {
		return fmt.Errorf("failed to make topics: %w", err)
	}
	var (
		logs   = make(chan types.Log, 128)
		config = ethereum.FilterQuery{
			Addresses: addresses,
			Topics:    topics,
		}
	)
	sub, err := c.c.SubscribeFilterLogs(ctx, config, logs)
	if err != nil {
		return fmt.Errorf("failed to start subscription: %w", err)
	}
	defer sub.Unsubscribe()

	// reuse
	var (
		event = new(ierc721.ContractTransfer)
		log   = types.Log{}
	)
	for {
		select {
		case <-ctx.Done():
			return nil
		case log = <-logs:
			if err = c.nftAbi.UnpackIntoInterface(event, "Transfer", log.Data); err != nil {
				return fmt.Errorf("failed to unpack log: %w", err)
			}
			event.Raw = log
			if err = callback(event); err != nil {
				return fmt.Errorf("failed to callback: %w", err)
			}
		case err = <-sub.Err():
			return fmt.Errorf("failed to watch transfer: %w", err)
		}
	}
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

func (c *EthHttpClient) getFactoryContract(ctx context.Context, address common.Address) (contract *factory.Factory, err error) {
	value, exist := c.cache.Get(address.Hex())

	if exist {
		contract, _ = value.(*factory.Factory)
		return
	}

	if contract, err = factory.NewFactory(address, c.c); err != nil {
		err = fmt.Errorf("failed to create contract: %w", err)
		return
	}
	c.cache.Add(address.Hex(), contract)
	return
}
