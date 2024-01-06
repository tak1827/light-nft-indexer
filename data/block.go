package data

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	PrefixBlock = []byte{byte('b')}
)

var _ StorableData = (*Block)(nil)

func NewBlock(height uint64, hash string, time uint64, t BlockType, sub string, now time.Time) (d *Block) {
	d = &Block{}
	d.Height = height
	d.Hash = hash
	d.Time = time
	d.Type = t
	d.SubIdentifier = sub
	d.CreatedAt = timestamppb.New(now)
	d.UpdatedAt = timestamppb.New(now)
	return
}

func (d *Block) Key() []byte {
	return append(PrefixBlock, []byte(fmt.Sprintf("%s%d%s%s", Separator, d.Type, Separator, d.SubIdentifier))...)
}

func (d *Block) Value() []byte {
	value, err := proto.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data: %w", err))
	}
	return value
}
