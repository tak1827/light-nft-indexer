package data

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	PrefixNFTContract     = []byte{byte('n')}
	PrefixToken           = []byte{byte('t')}
	PrefixTransferHistory = []byte{byte('h')}
)

var _ StorableData = (*NFTContract)(nil)

func NewNFTContract(address, name, symbol string, totalSupply uint64, tokenIds []string, now time.Time) (d *NFTContract) {
	d = &NFTContract{}
	d.Address = address
	d.Name = name
	d.Symbol = symbol
	d.TotalSupply = totalSupply
	d.TokenIds = tokenIds
	d.CreatedAt = timestamppb.New(now)
	d.UpdatedAt = timestamppb.New(now)
	return
}

func (d *NFTContract) Key() []byte {
	return append(PrefixNFTContract, []byte(fmt.Sprintf("%s%s", PrefixSeparator, d.Address))...)
}

func (d *NFTContract) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}

func (d *Token) Key() []byte {
	return append(PrefixToken, []byte(fmt.Sprintf("%s%s%s%s", PrefixSeparator, d.Address, PrefixSeparator, d.TokenId))...)
}

func (d *Token) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}

func NewTransferHistory(address, tokenId, from, to string, now *timestamppb.Timestamp) (d *TransferHistory) {
	d = &TransferHistory{}
	d.Address = address
	d.TokenId = tokenId
	d.From = from
	d.To = to
	d.CreatedAt = now
	return
}

// func (d *TransferHistory) Key() []byte {
// 	return append(PrefixNFTContract, []byte(fmt.Sprintf("%s%s%s%s", PrefixSeparator, d.Address, PrefixSeparator, d.TokenId))...)
// }

// func (d *TransferHistory) Value() []byte {
// 	value, err := proto.Marshal(d)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to marshal data: %w", err))
// 	}
// 	return value
// }
