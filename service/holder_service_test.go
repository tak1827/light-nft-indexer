package service

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/server"
	"github.com/tak1827/light-nft-indexer/store"
	"github.com/tak1827/light-nft-indexer/util/testutil"
)

var (
	owner1 = common.BigToAddress(big.NewInt(11))
	owner2 = common.BigToAddress(big.NewInt(12))
	owner3 = common.BigToAddress(big.NewInt(13))
)

func TestNewHolderService_ListHolderNftToken(t *testing.T) {
	var (
		dataSets = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex(), Name: "test1", Symbol: "T1"},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "101", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "102", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "103", Owner: owner2.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: twoAddress.Hex(), TokenId: "104", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
		}
		db      = testutil.InitDB()
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx)
		srv     = NewHolderService(&base)
	)
	defer db.Close()

	// no data
	_, err := srv.ListHolderNftToken(ctx, &ListHolderNftTokenRequest{ContractAddress: oneAddress.Hex(), WalletAddress: owner2.Hex()})
	require.ErrorIs(t, err, store.ErrNotFound)

	// with data
	batch := db.Batch()
	batch.Put(dataSets...)
	batch.Commit()
	res, err := srv.ListHolderNftToken(ctx, &ListHolderNftTokenRequest{ContractAddress: oneAddress.Hex(), WalletAddress: owner1.Hex()})
	require.NoError(t, err)
	require.Equal(t, "test1", res.GetNftContract().GetName())
	fmt.Println(res.GetTokens())
	require.Equal(t, 2, len(res.GetTokens()))
	var ids []string
	for _, token := range res.GetTokens() {
		ids = append(ids, token.GetTokenId())
	}
	require.EqualValues(t, []string{"101", "102"}, ids)
}

func TestNewHolderService_ListHolderAllNftToken(t *testing.T) {
	var (
		dataSets = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex(), Name: "test1", Symbol: "T1"},
			&data.NFTContract{Address: twoAddress.Hex(), Name: "test1", Symbol: "T1"},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "101", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "102", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: oneAddress.Hex(), TokenId: "103", Owner: owner2.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: twoAddress.Hex(), TokenId: "104", Owner: owner1.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
			&data.TokenOwnerIndex{Token: &data.Token{Address: twoAddress.Hex(), TokenId: "105", Owner: owner2.Hex(), Meta: &data.TokenMeta{Meta: "meta"}}},
		}
		db      = testutil.InitDB()
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx)
		srv     = NewHolderService(&base)
	)
	defer db.Close()

	// no data
	_, err := srv.ListHolderAllNftToken(ctx, &ListHolderAllNftTokenRequest{WalletAddress: owner3.Hex()})
	require.ErrorIs(t, err, store.ErrNotFound)

	// with data
	batch := db.Batch()
	batch.Put(dataSets...)
	batch.Commit()
	res, err := srv.ListHolderAllNftToken(ctx, &ListHolderAllNftTokenRequest{WalletAddress: owner1.Hex()})
	require.NoError(t, err)
	require.Equal(t, 2, len(res.GetNftContracts()))
	var ids []string
	for _, contract := range res.GetNftContracts() {
		for _, token := range contract.GetTokens() {
			ids = append(ids, token.GetTokenId())
		}
	}
	// sort ids
	sort.Strings(ids)
	require.EqualValues(t, []string{"101", "102", "104"}, ids)
}
