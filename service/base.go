package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/log"
)

type BaseService struct {
	Logger      zerolog.Logger
	ImageClient apiclient.ImageDownloadClient
}

func NewBaseService(ctx context.Context, image apiclient.ImageDownloadClient) (b BaseService, err error) {
	b.Logger = log.Service()
	b.ImageClient = image
	return
}
