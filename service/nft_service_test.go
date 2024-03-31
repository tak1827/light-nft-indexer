package service

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/server"
	"github.com/tak1827/light-nft-indexer/store"
	"github.com/tak1827/light-nft-indexer/util"
	"github.com/tak1827/light-nft-indexer/util/testutil"
)

var (
	oneAddress   = common.BigToAddress(big.NewInt(1))
	twoAddress   = common.BigToAddress(big.NewInt(2))
	threeAddress = common.BigToAddress(big.NewInt(3))
)

func TestNewNftService_GetNftContract(t *testing.T) {
	var (
		initalData = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex(), Name: "test1", Symbol: "T1"},
		}
		db      = testutil.InitDB(initalData...)
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx, nil)
		srv     = NewNftService(&base)
	)
	defer db.Close()

	tests := []struct {
		name    string
		req     *GetNftContractRequest
		tname   string
		tsymbol string
		err     error
	}{
		{
			name: "succeed",
			req: &GetNftContractRequest{
				ContractAddress: oneAddress.Hex(),
			},
			tname:   "test1",
			tsymbol: "T1",
		},
		{
			name: "not found",
			req: &GetNftContractRequest{
				ContractAddress: twoAddress.Hex(),
			},
			err: store.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.GetNftContract(ctx, tt.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tt.req.GetContractAddress(), res.GetNft().GetAddress())
			require.Equal(t, tt.tname, res.GetNft().GetName())
			require.Equal(t, tt.tsymbol, res.GetNft().GetSymbol())
		})
	}
}

func TestNewNftService_ListAllNftContract(t *testing.T) {
	var (
		dataSets = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex(), Name: "test1", Symbol: "T1"},
			&data.NFTContract{Address: twoAddress.Hex(), Name: "test2", Symbol: "T2"},
			&data.NFTContract{Address: threeAddress.Hex(), Name: "test3", Symbol: "T3"},
		}
		db      = testutil.InitDB()
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx, nil)
		srv     = NewNftService(&base)
	)
	defer db.Close()

	// no data
	res, err := srv.ListAllNftContract(ctx, &ListAllNftContractRequest{})
	require.ErrorIs(t, err, store.ErrNotFound)

	// with data
	batch := db.Batch()
	batch.Put(dataSets...)
	batch.Commit()
	res, err = srv.ListAllNftContract(ctx, &ListAllNftContractRequest{})
	require.NoError(t, err)
	require.Equal(t, 3, len(res.GetNfts()))
	var names []string
	for _, nft := range res.GetNfts() {
		names = append(names, nft.GetName())
	}
	require.EqualValues(t, []string{"test1", "test2", "test3"}, names)
}

func TestNewNftService_GetNftToken(t *testing.T) {
	var (
		tokenId    = "1243"
		initalData = []data.StorableData{
			&data.Token{Address: oneAddress.Hex(), TokenId: tokenId, Owner: "owner1", Meta: &data.TokenMeta{Origin: "meta"}},
			&data.TransferHistory{Address: oneAddress.Hex(), TokenId: tokenId, IndexLogInBlock: 1, From: "owner1", To: "owner2"},
			&data.TransferHistory{Address: oneAddress.Hex(), TokenId: tokenId, IndexLogInBlock: 2, From: "owner1", To: "owner2"},
		}
		db      = testutil.InitDB(initalData...)
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx, nil)
		srv     = NewNftService(&base)
	)
	defer db.Close()

	tests := []struct {
		name    string
		req     *GetNftTokenRequest
		tname   string
		tsymbol string
		err     error
	}{
		{
			name: "succeed",
			req: &GetNftTokenRequest{
				ContractAddress: oneAddress.Hex(),
				TokenId:         tokenId,
			},
			tname:   "test1",
			tsymbol: "T1",
		},
		{
			name: "invalid token id",
			req: &GetNftTokenRequest{
				ContractAddress: oneAddress.Hex(),
				TokenId:         "tokenid",
			},
			err: util.ErrNotNumber,
		},
		{
			name: "not found",
			req: &GetNftTokenRequest{
				ContractAddress: oneAddress.Hex(),
				TokenId:         "123",
			},
			err: store.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.GetNftToken(ctx, tt.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, "owner1", res.GetToken().GetOwner())
			require.Equal(t, 2, len(res.GetTransferHistories()))
		})
	}
}

func TestNewNftService_ListAllNftToken(t *testing.T) {
	var (
		dataSets = []data.StorableData{
			&data.Token{Address: oneAddress.Hex(), TokenId: "101", Owner: "owner1", Meta: &data.TokenMeta{Origin: "meta"}},
			&data.Token{Address: oneAddress.Hex(), TokenId: "102", Owner: "owner2", Meta: &data.TokenMeta{Origin: "meta"}},
			&data.Token{Address: oneAddress.Hex(), TokenId: "103", Owner: "owner3", Meta: &data.TokenMeta{Origin: "meta"}},
		}
		db      = testutil.InitDB()
		ctx     = server.CtxWithDB(context.Background(), db)
		base, _ = NewBaseService(ctx, nil)
		srv     = NewNftService(&base)
	)
	defer db.Close()

	// no data
	res, err := srv.ListAllNftToken(ctx, &ListAllNftTokenRequest{ContractAddress: oneAddress.Hex()})
	require.ErrorIs(t, err, store.ErrNotFound)

	// with data
	batch := db.Batch()
	batch.Put(dataSets...)
	batch.Commit()
	res, err = srv.ListAllNftToken(ctx, &ListAllNftTokenRequest{ContractAddress: oneAddress.Hex()})
	require.NoError(t, err)
	require.Equal(t, 3, len(res.GetTokens()))
	var ids []string
	for _, token := range res.GetTokens() {
		ids = append(ids, token.GetTokenId())
	}
	require.EqualValues(t, []string{"101", "102", "103"}, ids)
}
