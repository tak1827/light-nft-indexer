package watchlog

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/store"
)

var (
	zeroAddress  = common.BigToAddress(big.NewInt(0))
	oneAddress   = common.BigToAddress(big.NewInt(1))
	twoAddress   = common.BigToAddress(big.NewInt(2))
	threeAddress = common.BigToAddress(big.NewInt(3))

	acc1 = common.BigToAddress(big.NewInt(1001))
	acc2 = common.BigToAddress(big.NewInt(1002))
	acc3 = common.BigToAddress(big.NewInt(1003))
)

func initDB(ds ...data.StorableData) store.DB {
	var (
		db, _ = store.NewPebbleDB("", true, nil)
		batch = db.Batch()
	)
	defer batch.Close()

	batch.Put(ds...)
	batch.Commit()

	return &db
}

func TestWatcher(t *testing.T) {
	var (
		initalData = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex()},
			&data.Token{Address: oneAddress.Hex(), TokenId: "1"},
		}
		db = initDB(initalData...)
		fe = []*factory.FactoryNFTCreated{{Nft: twoAddress}, {Nft: threeAddress}}
		te = []*ierc721.ContractTransfer{
			{
				TokenId: big.NewInt(1),
				To:      acc1,
				Raw: types.Log{
					Address: oneAddress,
				},
			},
			{
				TokenId: big.NewInt(1),
				To:      acc2,
				Raw: types.Log{
					Address: twoAddress,
				},
			},
			{
				TokenId: big.NewInt(2),
				To:      acc3,
				Raw: types.Log{
					Address: oneAddress,
				},
			},
		}
		client = apiclient.NewMockChainClient(fe, te, 10)
		w      = NewWatcher(db, &client, 100)
		err    error
	)
	defer db.Close()

	db.Batch()

	ctx, cancel := context.WithCancel(context.Background())

	err = w.Start(ctx, common.Address{})
	require.NoError(t, err)

	time.Sleep(150 * time.Millisecond)
	cancel()

	require.NoError(t, w.Close())
	require.NoError(t, <-w.Err())

	// check db
	nfts := []*data.NFTContract{{}}
	require.NoError(t, db.List(data.PrefixNFTContract, &nfts))
	require.Len(t, nfts, 3)
	require.Equal(t, fe[0].Nft.Hex(), nfts[1].Address)
	require.Equal(t, fe[1].Nft.Hex(), nfts[2].Address)

	tokens := []*data.Token{{}}
	require.NoError(t, db.List(data.PrefixToken, &tokens))
	sort.Slice(tokens, func(i, j int) bool { return tokens[i].Address < tokens[j].Address })
	require.Len(t, tokens, 3)
	require.Equal(t, te[0].Raw.Address.Hex(), tokens[0].Address)
	require.Equal(t, te[0].TokenId.String(), tokens[0].TokenId)
	require.Equal(t, te[0].To.Hex(), tokens[0].Owner)
	require.Equal(t, te[1].Raw.Address.Hex(), tokens[2].Address)
	require.Equal(t, te[1].TokenId.String(), tokens[2].TokenId)
	require.Equal(t, te[1].To.Hex(), tokens[2].Owner)
	require.Equal(t, te[2].Raw.Address.Hex(), tokens[1].Address)
	require.Equal(t, te[2].TokenId.String(), tokens[1].TokenId)
	require.Equal(t, te[2].To.Hex(), tokens[1].Owner)
}

func TestErrHand(t *testing.T) {
	var (
		db, _       = store.NewPebbleDB("", true, nil)
		client      = apiclient.NewMockChainClient([]*factory.FactoryNFTCreated{}, []*ierc721.ContractTransfer{}, 10)
		w           = NewWatcher(&db, &client, 100)
		expectedErr = errors.New("test")
		err         error
	)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, w.Start(ctx, common.Address{}))

	time.Sleep(10 * time.Millisecond)
	client.EmitErr(expectedErr)

	err = <-w.Err()
	require.True(t, errors.Is(err, expectedErr))

	require.NoError(t, w.Close())
}
