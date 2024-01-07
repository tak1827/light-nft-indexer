package apiclient

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
)

var (
	ErrConnectRefused = errors.New("connection endppint refused")
)

type ChainHttpClient interface {
	FetchTransferLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*ierc721.ContractTransfer, nextStart uint64, err error)
	FetchFactoryLog(ctx context.Context, address common.Address, startHeight uint64, endHeight *uint64) (events []*factory.FactoryNFTCreated, nextStart uint64, err error)

	WatchTransferLog(ctx context.Context, addresses []common.Address, callback func(*ierc721.ContractTransfer) error) error
	WatchFactoryLog(ctx context.Context, address common.Address, callback func(*factory.FactoryNFTCreated) error) error

	FetchNFTInfo(ctx context.Context, d *data.NFTContract) (err error)
	GetTokenMeta(ctx context.Context, d *data.Token) (err error)
}
