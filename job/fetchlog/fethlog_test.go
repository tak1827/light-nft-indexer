package fetchlog

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/contract/factory"
	"github.com/tak1827/light-nft-indexer/contract/ierc721"
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/util/testutil"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestFetcher_FetchNFTs(t *testing.T) {
	var (
		fe         = []*factory.FactoryNFTCreated{{Nft: oneAddress}, {Nft: twoAddress}}
		initalData = []data.StorableData{
			&data.NFTContract{Address: oneAddress.Hex()},
			&data.NFTContract{Address: twoAddress.Hex()},
			&data.Block{Type: data.BlockType_BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED, Height: 10},
		}
		endHeight = uint64(100)
		db        = testutil.InitDB(initalData...)
		client    = apiclient.NewMockChainClient(fe, nil, &endHeight, 10)
		f         = NewFetcher(db, &client)
		ctx       = context.Background()
		now       = timestamppb.Now()
	)
	defer db.Close()

	_, err := f.FetchNFTs(ctx, zeroAddress, now, true)
	require.NoError(t, err)

	// checking the nfts are fetched
	nfts := []*data.NFTContract{{}}
	err = db.List(data.PrefixNFTContract, &nfts)
	require.NoError(t, err)
	require.Len(t, nfts, 2)
	require.Equal(t, oneAddress.Hex(), nfts[0].Address)
	require.Equal(t, twoAddress.Hex(), nfts[1].Address)

	// checking the height of next fetch block is updated
	b, err := f.getBlock(data.BlockType_BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED, "", now)
	require.NoError(t, err)
	require.Equal(t, endHeight, b.GetHeight())
}

func TestFetcher_FetchTokens(t *testing.T) {
	var (
		nft = data.NFTContract{Address: oneAddress.Hex()}
		te  = []*ierc721.ContractTransfer{
			{From: acc1, To: acc2, TokenId: big.NewInt(101)},
			{From: acc2, To: acc3, TokenId: big.NewInt(102)},
		}
		initalData = []data.StorableData{
			&nft,
			&data.Block{Type: data.BlockType_BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED, Height: 10},
		}
		endHeight = uint64(100)
		db        = testutil.InitDB(initalData...)
		client    = apiclient.NewMockChainClient(nil, te, &endHeight, 10)
		f         = NewFetcher(db, &client)
		ctx       = context.Background()
		now       = timestamppb.Now()
	)
	defer db.Close()

	_, err := f.FetchTokens(ctx, zeroAddress, now, &nft, true)
	require.NoError(t, err)

	// checking the tokens are fetched
	tokens := []*data.Token{{}}
	err = db.List(append(data.PrefixToken, []byte(fmt.Sprintf("%s%s", data.Separator, nft.Address))...), &tokens)
	require.NoError(t, err)
	require.Len(t, tokens, 2)
	// sort the tokens by id
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].TokenId < tokens[j].TokenId
	})
	require.Equal(t, acc2.Hex(), tokens[0].GetOwner())
	require.Equal(t, acc3.Hex(), tokens[1].GetOwner())
	require.Equal(t, "meta", tokens[0].GetMeta().GetMeta())

	// checking the height of next fetch block is updated
	b, err := f.getBlock(data.BlockType_BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED, nft.Address, now)
	require.NoError(t, err)
	require.Equal(t, endHeight, b.GetHeight())
}
