package service

import (
	"context"
	"fmt"

	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/server"
)

type NftService struct {
	*BaseService
}

func NewNftService(base *BaseService) NftServer {
	return &NftService{
		BaseService: base,
	}
}

func (s *NftService) GetNftContract(ctx context.Context, req *GetNftContractRequest) (*GetNftContractResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("at GetNftContract.ValidateBasic: %w", err)
	}

	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	nft := data.NFTContract{Address: req.GetContractAddress()}
	if err := db.Get(nft.Key(), &nft); err != nil {
		return nil, fmt.Errorf("failed to get nft contract: %w", err)
	}

	return &GetNftContractResponse{Nft: &nft}, nil
}

func (s *NftService) ListAllNftContract(ctx context.Context, req *ListAllNftContractRequest) (*ListAllNftContractResponse, error) {
	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	nfts := []*data.NFTContract{{}}
	if err := db.List(data.PrefixNFTContract, &nfts); err != nil {
		return nil, fmt.Errorf("failed to get list of nft contracts: %w", err)
	}

	return &ListAllNftContractResponse{Nfts: nfts}, nil
}

func (s *NftService) GetNftToken(ctx context.Context, req *GetNftTokenRequest) (*GetNftTokenResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("at GetNftToken.ValidateBasic: %w", err)
	}

	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	token := data.Token{Address: req.GetContractAddress(), TokenId: req.GetTokenId()}
	if err := db.Get(token.Key(), &token); err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	histories := []*data.TransferHistory{{}}
	if err := db.List(data.PrefixTransferHistoryByTokenId(req.GetContractAddress(), req.GetTokenId()), &histories); err != nil {
		return nil, fmt.Errorf("failed to get list of nft contracts: %w", err)
	}

	return &GetNftTokenResponse{Token: &token, TransferHistories: histories}, nil
}

func (s *NftService) ListAllNftToken(ctx context.Context, req *ListAllNftTokenRequest) (*ListAllNftTokenResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("at GetNftToken.ValidateBasic: %w", err)
	}

	db, err := server.DBFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("at server.DBFromCtx: %w", err)
	}

	tokens := []*data.Token{{}}
	if err := db.List(data.PrefixToken, &tokens); err != nil {
		return nil, fmt.Errorf("failed to get list of nft contracts: %w", err)
	}

	return &ListAllNftTokenResponse{
		Tokens: tokenToTokenMinis(tokens),
	}, nil
}
