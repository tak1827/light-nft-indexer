package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/tak1827/light-nft-indexer/log"
)

type BaseService struct {
	Logger zerolog.Logger
}

func NewBaseService(ctx context.Context) (b BaseService, err error) {
	b.Logger = log.Service()
	return
}
