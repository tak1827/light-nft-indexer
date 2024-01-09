package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/data"
	"google.golang.org/protobuf/proto"
)

func TestPutUpdateDeleteGet(t *testing.T) {
	var (
		db, _    = NewPebbleDB("", true, nil)
		testdata = data.NewNFTContract("addr.1", "name.1", "symbol.1", 1, []string{"tkn.1", "tkn.2"}, time.Now())
		err      error
	)
	defer db.Close()

	// put
	batch := db.Batch()
	batch.Put(testdata)
	err = batch.Commit()
	require.NoError(t, err)
	batch.Close()

	// get
	var result data.NFTContract
	err = db.Get(testdata.Key(), &result)
	require.NoError(t, err)
	require.Equal(t, testdata.Value(), result.Value())

	// update by put
	db.cache.Clear()
	batch = db.Batch()
	testdata.Name = "name.2"
	batch.Put(testdata)
	err = batch.Commit()
	require.NoError(t, err)
	var result2 data.NFTContract
	err = db.Get(testdata.Key(), &result2)
	require.NoError(t, err)
	require.Equal(t, testdata.Value(), result2.Value())
	batch.Close()

	// delete
	batch = db.Batch()
	batch.Delete(testdata)
	err = batch.Commit()
	require.NoError(t, err)
	err = db.Get(testdata.Key(), &result)
	require.ErrorAs(t, err, &ErrNotFound)
	batch.Close()
}

func TestList(t *testing.T) {
	var (
		db, _    = NewPebbleDB("", true, nil)
		testdata = []*data.NFTContract{
			data.NewNFTContract("addr.1", "name.1", "symbol.1", 1, []string{"tkn.1", "tkn.2"}, time.Now()),
			data.NewNFTContract("addr.2", "name.2", "symbol.2", 2, []string{"tkn.3", "tkn.4"}, time.Now()),
			data.NewNFTContract("addr.3", "name.3", "symbol.3", 3, []string{"tkn.5", "tkn.6"}, time.Now()),
		}
		batch = db.Batch()
	)
	defer db.Close()
	defer batch.Close()

	batch.Put(testdata[0], testdata[1], testdata[2])
	batch.Commit()

	results := []*data.NFTContract{{}}
	err := db.List(data.PrefixNFTContract, &results)
	require.NoError(t, err)
	require.Equal(t, len(testdata), len(results))
	for i := range results {
		require.Equal(t, testdata[i].Value(), results[i].Value())
	}
}

func TestDeleteAll(t *testing.T) {
	var (
		db, _    = NewPebbleDB("", true, nil)
		testdata = []*data.NFTContract{
			data.NewNFTContract("addr.1", "name.1", "symbol.1", 1, []string{"tkn.1", "tkn.2"}, time.Now()),
			data.NewNFTContract("addr.2", "name.2", "symbol.2", 2, []string{"tkn.3", "tkn.4"}, time.Now()),
			data.NewNFTContract("addr.3", "name.3", "symbol.3", 3, []string{"tkn.5", "tkn.6"}, time.Now()),
		}
		batch = db.Batch()
		// keys  = [][]byte{{0x00}, {0x00, 0x00}, {0x00, 0x01}, {0xff}, {0xff, 0xff}, {0xff, 0xfe}}
	)
	defer db.Close()
	defer batch.Close()

	batch.Put(testdata[0], testdata[1], testdata[2])
	// for i := range keys {
	// 	batch.(*PebbleBatch).pb.Set(keys[i], keys[i], pebble.Sync)
	// }
	batch.Commit()

	err := db.DeleteAll(data.PrefixNFTContract)
	require.NoError(t, err)

	var result data.NFTContract
	for i := range testdata {
		err = db.Get(testdata[i].Key(), &result)
		require.ErrorAs(t, err, &ErrNotFound)
	}
}

func TestBatchEmptyLenContents(t *testing.T) {
	var (
		db, _  = NewPebbleDB("", true, nil)
		conunt = 0
		batch  = db.Batch()
	)
	defer db.Close()
	defer batch.Close()

	for conunt < 2 {
		testdata := []*data.NFTContract{
			data.NewNFTContract("addr.1", "name.1", "symbol.1", 1, []string{"tkn.1", "tkn.2"}, time.Now()),
			data.NewNFTContract("addr.2", "name.2", "symbol.2", 2, []string{"tkn.3", "tkn.4"}, time.Now()),
			data.NewNFTContract("addr.3", "name.3", "symbol.3", 3, []string{"tkn.5", "tkn.6"}, time.Now()),
		}

		// test length
		batch.Put(testdata[0], testdata[1], testdata[2])
		require.Equal(t, 3, batch.Len())

		// test contents
		_, values := batch.Contents()
		dests := make([]proto.Message, len(values))
		for i := range dests {
			dests[i] = &data.NFTContract{}
		}
		Unmarshal(values, dests)
		for i := range dests {
			require.Equal(t, testdata[i].Address, dests[i].(*data.NFTContract).Address)
		}

		// test empty
		batch.Reset()
		require.True(t, batch.Empty())
		require.Equal(t, 0, batch.Len())
		k, v := batch.Contents()
		require.Equal(t, [][]byte{}, k)
		require.Equal(t, [][]byte{}, v)

		conunt++
	}
}
