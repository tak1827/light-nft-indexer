package tools

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/data"
)

func Test_FillTokenMeta(t *testing.T) {
	var (
		nftName     = "DaveStarbelly"
		nativeImage = "native-image"
		filename    = "3.png"
		nftURL      = "https://storage.googleapis.com/opensea-prod.appspot.com/puffs/" + filename
		ctx         = context.Background()

		client = apiclient.NewMockImageClient(map[string]string{
			nftURL: "3.png",
		})

		now = time.Now().Unix()
	)

	tests := []struct {
		name      string
		meta      *data.TokenMeta
		now       int64
		imageType data.ImageType
		imageData string
		err       string
	}{
		{
			name: "native image",
			meta: &data.TokenMeta{
				Origin: `{"description":"Friendly OpenSea Creature that enjoys long swims in the ocean.","external_url":"https://openseacreatures.io/3","image":"` + nativeImage + `","name":"` + nftName + `","attributes":["aaa","bbb"]}`,
			},
			imageType: data.ImageType_IMAGE_TYPE_NATIVE,
			imageData: nativeImage,
			now:       now,
		},
		{
			name: "remote image",
			meta: &data.TokenMeta{
				Origin: `{"description":"Friendly OpenSea Creature that enjoys long swims in the ocean.","external_url":"https://openseacreatures.io/3","image":"` + nftURL + `","name":"` + nftName + `","attributes":["aaa","bbb"]}`,
			},
			imageType: data.ImageType_IMAGE_TYPE_REFERENCE,
			imageData: apiclient.DefaultBaseImageURL + "/" + nftName + "/" + filename,
			now:       now,
		},
		{
			name: "no origin",
			meta: &data.TokenMeta{},
			now:  now,
			err:  "origin is empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FillTokenMeta(ctx, &client, tt.meta, tt.now)
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, uint32(now)+TokenURLValidDuration, tt.meta.TimeOfExpire)
			require.Equal(t, nftName, tt.meta.Standard.Name)
			require.Equal(t, tt.imageData, tt.meta.Image.Data)
			require.Equal(t, tt.imageType, tt.meta.Image.Type)
		})
	}
}

func Test_IsExpiredTokenMeta(t *testing.T) {
	var (
		ctx  = context.Background()
		now  = time.Now().Unix()
		meta = &data.TokenMeta{
			TimeOfExpire: uint32(now) - 1,
			Image:        &data.TokenMetaImage{},
		}
	)
	expires := IsExpiredTokenMeta(ctx, meta, now)
	require.True(t, expires)
}
