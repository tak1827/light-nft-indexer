package data

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	PrefixBlock = []byte("block")
)

var _ StorableData = (*Block)(nil)

func NewBlock(height uint64, hash string, time uint64, t BlockType, now time.Time) (d *Block) {
	d = &Block{}
	d.Height = height
	d.Hash = hash
	d.Time = time
	d.Type = t
	d.CreatedAt = timestamppb.New(now)
	d.UpdatedAt = timestamppb.New(now)
	return
}

func (d *Block) Key() []byte {
	return append(append(PrefixBlock, PrefixSeparator), []byte(fmt.Sprintf("%d", d.Type))...)
}

func (d *Block) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}
