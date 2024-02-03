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
	PrefixTokenOwnerIndex = []byte{byte('o')}
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
	return append(PrefixNFTContract, []byte(fmt.Sprintf("%s%s", Separator, d.Address))...)
}

func (d *NFTContract) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}

func (d *Token) Key() []byte {
	return append(PrefixToken, []byte(fmt.Sprintf("%s%s%s%s", Separator, d.Address, Separator, d.TokenId))...)
}

func (d *Token) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}

type TokenOwnerIndex struct {
	*Token
}

func (d *TokenOwnerIndex) Key() []byte {
	return append(PrefixTokenOwnerIndex, []byte(fmt.Sprintf("%s%s%s%s", Separator, d.Owner, Separator, d.TokenId))...)
}

func (d *TokenOwnerIndex) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}

func NewTokenOwnerIndex(d *Token) *TokenOwnerIndex {
	return &TokenOwnerIndex{Token: d}
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

func PrefixTransferHistoryByTokenId(address, tokenId string) []byte {
	return append(PrefixTransferHistory, []byte(fmt.Sprintf("%s%s%s%s", Separator, address, Separator, tokenId))...)
}

func (d *TransferHistory) Key() []byte {
	return append(PrefixTransferHistoryByTokenId(d.Address, d.TokenId), []byte(fmt.Sprintf("%s%d%s%d", Separator, d.BlockNumber, Separator, d.IndexLogInBlock))...)
}

func (d *TransferHistory) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}
