package store

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/tak1827/go-cache/lru"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Batch interface {
	Put(item ...data.StorableData)
	PutWithTime(t *timestamppb.Timestamp, item ...data.StorableData)
	Delete(item data.StorableData)

	Commit() error

	Reset()
	Close() error

	Empty() bool
	Len() int
	Contents() ([][]byte, [][]byte)
}

var _ Batch = (*PebbleBatch)(nil)

type PebbleBatch struct {
	pb    *pebble.Batch
	cache *lru.LRUCache
}

func NewBatch(db *pebble.DB, c *lru.LRUCache) Batch {
	return &PebbleBatch{
		pb:    db.NewBatch(),
		cache: c,
	}
}

func (b *PebbleBatch) Put(item ...data.StorableData) {
	for i := range item {
		b.pb.Set(item[i].Key(), item[i].Value(), pebble.Sync)
		clearCache(b.cache, item[i].Key())
	}
}

func (b *PebbleBatch) PutWithTime(t *timestamppb.Timestamp, item ...data.StorableData) {
	for i := range item {
		util.SetValueByName("UpdatedAt", t, item[i])
		if item[i].GetCreatedAt() == nil {
			util.SetValueByName("CreatedAt", t, item[i])
		}
	}
	b.Put(item...)
}

func (b *PebbleBatch) Delete(value data.StorableData) {
	b.pb.Delete(value.Key(), pebble.Sync)
	clearCache(b.cache, value.Key())
}

// func (b *PebbleBatch) Update(value data.StorableData) {
// 	b.Items = append(b.Items, value)
// 	b.pb.Merge(value.Key(), value.Value(), pebble.Sync)
// 	clearCache(b.cache, value.Key())
// }

func (b *PebbleBatch) Commit() error {
	if err := b.pb.Commit(pebble.Sync); err != nil {
		return fmt.Errorf("failed to commit batch: %w", err)
	}
	b.pb.Reset()
	return nil
}

func (b *PebbleBatch) Reset() {
	b.pb.Reset()
}

func (b *PebbleBatch) Close() error {
	return b.pb.Close()
}

func (b *PebbleBatch) Empty() bool {
	return b.pb.Empty()
}

func (b *PebbleBatch) Len() int {
	return int(b.pb.Count())
}

func (b *PebbleBatch) Contents() ([][]byte, [][]byte) {
	var (
		reader = b.pb.Reader()
		keys   = make([][]byte, b.pb.Count())
		values = make([][]byte, b.pb.Count())
	)
	for i := uint32(0); i < b.pb.Count(); i++ {
		_, key, value, ok := reader.Next()
		if !ok {
			panic(fmt.Sprintf("inconsistent batch. expected %d entries, found %d", b.pb.Count(), i-1))
		}
		keys[i] = key
		values[i] = value
	}

	return keys, values
}

func clearCache(c *lru.LRUCache, key []byte) {
	c.Remove(string(key))
}

func Unmarshal(values [][]byte, dests []proto.Message) (err error) {
	// handle unexpected panic
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("invalid dests: %v", rec)
		}
	}()

	if len(values) != len(dests) {
		return fmt.Errorf("invalid dests: length of values and dests must be the same")
	}

	for i := range values {
		if err = proto.Unmarshal(values[i], dests[i]); err != nil {
			return fmt.Errorf("failed to unmarshal %dth index of value: %w", i+1, err)
		}
	}
	return
}
