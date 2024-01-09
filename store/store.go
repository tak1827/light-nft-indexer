package store

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/tak1827/go-cache/lru"
	"github.com/tak1827/light-nft-indexer/util"
	"google.golang.org/protobuf/proto"
)

const (
	DefaultCacheSize = 1024
	DefaultCacheTTL  = 60 * 3 // 3min
)

var (
	ErrNotFound         = errors.New("data not found")
	ErrInvalidPagingKey = errors.New("invalid paging key")

	keyUpperBound = func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil // no upper-bound
	}

	prefixIterOptions = func(prefix []byte) *pebble.IterOptions {
		return &pebble.IterOptions{
			LowerBound: prefix,
			UpperBound: keyUpperBound(prefix),
		}
	}
)

type DB interface {
	Batch() Batch
	Get(key []byte, result proto.Message) (err error)
	List(prefix []byte, results interface{}) (err error)
	DeleteAll(prefix []byte) error
	Clear() error
	Close() error
}

var _ DB = (*PebbleDB)(nil)

type PebbleDB struct {
	db   *pebble.DB
	path string

	cache     lru.LRUCache
	cacheSize int
	cacheTTL  int
}

func NewPebbleDB(path string, isMem bool, opt *pebble.Options) (d PebbleDB, err error) {
	if opt == nil {
		opt = &pebble.Options{}
	}
	if isMem {
		opt.FS = vfs.NewMem()
	}
	// if opt.Comparer == nil {
	// 	opt.Comparer = pebble.DefaultComparer
	// 	opt.Comparer.Split = func(a []byte) int { return len(a) }
	// }
	db, err := pebble.Open(path, opt)
	if err != nil {
		err = fmt.Errorf("failed to open pebble db: %w", err)
		return
	}

	d.cacheSize = DefaultCacheSize
	d.cacheTTL = DefaultCacheTTL
	d.db = db
	d.cache = lru.NewCache(d.cacheSize, d.cacheTTL)
	return
}

func (d *PebbleDB) Close() error {
	d.cache.Clear()
	return d.db.Close()
}

func (d *PebbleDB) Batch() Batch {
	return NewBatch(d.db, &d.cache)
}

func (d *PebbleDB) Get(key []byte, result proto.Message) (err error) {
	value, exist := d.getCache(key)

	if !exist {
		var closer io.Closer
		value, closer, err = d.db.Get(key)
		if closer != nil {
			closer.Close()
		}
		if err != nil {
			if errors.Is(err, pebble.ErrNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("failed to get value from pebble db: %w", err)
		}
	}

	if err := proto.Unmarshal(value, result); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	d.addCache(key, value)
	return nil
}

func (d *PebbleDB) List(prefix []byte, results interface{}) (err error) {
	if !util.IsPointer(results) {
		return errors.New("results should be pointer")
	}

	var (
		origin   = reflect.ValueOf(results)
		vResults = origin.Elem()
	)
	if vResults.Kind() != reflect.Slice {
		return errors.New("results should be slice")
	}
	if vResults.Len() == 0 {
		return ErrNotFound
	}
	protoType := reflect.TypeOf((*proto.Message)(nil)).Elem()
	if !vResults.Index(0).Type().Implements(protoType) {
		return fmt.Errorf("results array must contain proto.Message")
	}

	var (
		iter      = d.db.NewIter(prefixIterOptions(prefix))
		i         = 0
		firstElem = vResults.Index(0)
	)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}
	for iter.First(); iter.Valid(); iter.Next() {
		// if the results array is smaller than the number of results in the db, we need to append
		if vResults.Len() <= i {
			vResults = reflect.Append(vResults, reflect.New(firstElem.Type()))
		}

		elem := vResults.Index(i)
		if err = proto.Unmarshal(iter.Value(), elem.Interface().(proto.Message)); err != nil {
			return fmt.Errorf("failed to unmarshal %dth index of value: %w", i+1, err)
		}
		i++
	}
	if err = iter.Error(); err != nil {
		return fmt.Errorf("failed to iterate: %w", err)
	}
	if err = iter.Close(); err != nil {
		return fmt.Errorf("failed to close iterator: %w", err)
	}

	origin.Elem().Set(vResults)
	return
}

func (d *PebbleDB) DeleteAll(prefix []byte) error {
	if err := d.db.DeleteRange(prefix, keyUpperBound(prefix), pebble.Sync); err != nil {
		return fmt.Errorf("failed to delete all keys: %w", err)
	}
	return nil
}

// NOTE: don't work as expected
func (d *PebbleDB) Clear() error {
	if err := d.db.RangeKeyDelete([]byte{0x00}, []byte{0xff}, pebble.Sync); err != nil {
		return fmt.Errorf("failed to delete all keys: %w", err)
	}
	return nil
}

func (d *PebbleDB) getCache(key []byte) (result []byte, exist bool) {
	value, exist := d.cache.Get(string(key))
	if !exist {
		return
	}
	result = value.([]byte)
	return
}

func (d *PebbleDB) addCache(key []byte, result interface{}) {
	d.cache.Add(string(key), result)
}

func (d *PebbleDB) removeCache(key []byte) {
	d.cache.Remove(string(key))
}

func (d *PebbleDB) clearCache() {
	d.cache.Clear()
}
