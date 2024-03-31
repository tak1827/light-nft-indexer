package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/tak1827/light-nft-indexer/apiclient"
	"github.com/tak1827/light-nft-indexer/data"
)

const (
	TokenURLValidDuration = 60 * 60 * 24 * 7 // 7days
)

func IsExpiredTokenMeta(ctx context.Context, meta *data.TokenMeta, now int64) bool {
	if meta == nil {
		return true
	}
	if now == 0 {
		now = time.Now().Unix()
	}
	if meta.Image != nil {
		return uint32(now) >= meta.TimeOfExpire
	}
	return true
}

func FillTokenMeta(ctx context.Context, imgclient apiclient.ImageDownloadClient, meta *data.TokenMeta, now int64) (err error) {
	if meta == nil || meta.Origin == "" {
		return errors.New("meta is nil or origin is empty")
	}
	if now == 0 {
		now = time.Now().Unix()
	}

	// update expire time
	meta.TimeOfExpire = uint32(now) + TokenURLValidDuration

	// fill meta standard
	if meta.Standard == nil {
		meta.Standard = &data.MetaStandard{}
	}
	if err = json.Unmarshal([]byte(meta.Origin), meta.Standard); err != nil {
		return fmt.Errorf("failed to unmarshal meta: %w", err)
	}

	// fill meta image
	if meta.Image == nil {
		meta.Image = &data.TokenMetaImage{}
	}
	if meta.Standard.Image != "" {
		if !isHTTPURL(meta.Standard.Image) {
			// native image
			meta.Image.Type = data.ImageType_IMAGE_TYPE_NATIVE
			meta.Image.Data = meta.Standard.Image
		} else {
			// remote image
			meta.Image.Type = data.ImageType_IMAGE_TYPE_REFERENCE
			tag := ""
			if meta.Standard.Name != "" {
				tag = meta.Standard.Name
			}
			if meta.Image.Data, err = imgclient.Download(ctx, meta.Standard.Image, tag, true); err != nil {
				return fmt.Errorf("failed to download image: %w", err)
			}
		}
	}

	return nil
}

// isHTTPURL checks if the given string is a valid HTTP URL.
func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
