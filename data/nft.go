package data

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	PrefixNFTContract = []byte("nftcontract")
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
	return append(append(PrefixNFTContract, PrefixSeparator), []byte(d.Address)...)
}

func (d *NFTContract) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}
