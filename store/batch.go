package store

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/tak1827/go-cache/lru"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Batch interface {
	Put(item ...data.StorableData)
	PutWithTime(t *timestamppb.Timestamp, item ...data.StorableData)
	Delete(item data.StorableData)
	Commit() error
	Reset()
	Close() error
	Contents() []data.StorableData
}

var _ Batch = (*PebbleBatch)(nil)

type PebbleBatch struct {
	pb    *pebble.Batch
	cache *lru.LRUCache
	Items []data.StorableData
}

func NewBatch(db *pebble.DB, c *lru.LRUCache) Batch {
	return &PebbleBatch{
		pb:    db.NewBatch(),
		cache: c,
		Items: []data.StorableData{},
	}
}

func (b *PebbleBatch) Put(item ...data.StorableData) {
	for i := range item {
		b.pb.Set(item[i].Key(), item[i].Value(), pebble.Sync)
	}
	b.Items = append(b.Items, item...)
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
	return nil
}

func (b *PebbleBatch) Reset() {
	b.pb.Reset()
}

func (b *PebbleBatch) Close() error {
	return b.pb.Close()
}

func (b *PebbleBatch) Contents() []data.StorableData {
	return b.Items
}

func clearCache(c *lru.LRUCache, key []byte) {
	c.Remove(string(key))
}
