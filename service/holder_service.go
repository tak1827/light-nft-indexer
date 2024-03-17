package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/server"
)

type HolderService struct {
	*BaseService
}

func NewHolderService(base *BaseService) HolderServer {
	return &HolderService{
		BaseService: base,
	}
}

func (s *HolderService) ListHolderNftToken(ctx context.Context, req *ListHolderNftTokenRequest) (*ListHolderNftTokenResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("at ListHolderNftToken.ValidateBasic: %w", err)
	}

	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	nft := data.NFTContract{Address: req.GetContractAddress()}
	if err := db.Get(nft.Key(), &nft); err != nil {
		return nil, fmt.Errorf("failed to get nft contract: %w", err)
	}

	tokens := []*data.Token{{}}
	prefix := append(data.PrefixTokenOwnerIndex, []byte(fmt.Sprintf("%s%s%s%s", data.Separator, req.GetWalletAddress(), data.Separator, req.GetContractAddress()))...)
	if err := db.List(prefix, &tokens); err != nil {
		return nil, fmt.Errorf("failed to get list of nft tokens: %w", err)
	}

	return &ListHolderNftTokenResponse{NftContract: &nft, Tokens: tokenToTokenMinis(tokens)}, nil
}

func (s *HolderService) ListHolderAllNftToken(ctx context.Context, req *ListHolderAllNftTokenRequest) (*ListHolderAllNftTokenResponse, error) {
	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	tokens := []*data.Token{{}}
	prefix := append(data.PrefixTokenOwnerIndex, []byte(fmt.Sprintf("%s%s", data.Separator, req.GetWalletAddress()))...)
	if err := db.List(prefix, &tokens); err != nil {
		return nil, fmt.Errorf("failed to get list of nft tokens: %w", err)
	}

	// sort by contract address
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Address < tokens[j].Address
	})

	contracts := []*ContractWithToken{}
	for _, token := range tokens {
		if len(contracts) == 0 || contracts[len(contracts)-1].NftContract.Address != token.Address {
			nft := data.NFTContract{Address: token.Address}
			if err := db.Get(nft.Key(), &nft); err != nil {
				return nil, fmt.Errorf("failed to get nft contract: %w", err)
			}
			mini := tokenToTokenMini(token)
			contracts = append(contracts, &ContractWithToken{NftContract: &nft, Tokens: []*TokenMini{mini}})
		} else {
			contracts[len(contracts)-1].Tokens = append(contracts[len(contracts)-1].Tokens, tokenToTokenMini(token))
		}
	}

	return &ListHolderAllNftTokenResponse{NftContracts: contracts}, nil
}

func tokenToTokenMini(token *data.Token) *TokenMini {
	return &TokenMini{
		TokenId: token.TokenId,
		Owner:   token.Owner,
		Meta:    token.Meta,
	}
}

func tokenToTokenMinis(tokens []*data.Token) []*TokenMini {
	tokenMinis := make([]*TokenMini, len(tokens))
	for i, token := range tokens {
		tokenMinis[i] = tokenToTokenMini(token)
	}
	return tokenMinis
}
