package apiclient

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestTrustdocEndpoint = "https://api.test.trustdock.io/v2/"
	TestTrustdocApiToken = "nfxTUWaXHggW7sn4aWXpUVGR"
	TestPlanId1          = "eaa8bcc3-6ec9-4c25-b4ab-54447e1cd4cb"
	TestPlanId2          = "f18a07e7-b667-47fe-a146-c390df76d241"
)

func TestDownload(t *testing.T) {
	var (
		ctx          = context.Background()
		baseLocation = "./static-test"
		baseImageURL = "https://example.com"
		filename     = "logotop.svg"
		imageURL     = "https://bitcoin.org/img/icons/" + filename
		tag          = "imagetest"
		expectedURL  = baseImageURL + "/" + tag + "/" + filename
	)
	iclient, err := NewLocalImageDownloadClient(ctx, baseLocation, baseImageURL)
	require.NoError(t, err)

	// overwrite
	url, err := iclient.Download(ctx, imageURL, tag, true)
	require.NoError(t, err)
	require.Equal(t, expectedURL, url)

	// not overwrite
	_, err = iclient.Download(ctx, imageURL, tag, false)
	require.ErrorIs(t, err, ErrImageAlreadyExists)

	// clear
	os.RemoveAll(baseLocation)
}
