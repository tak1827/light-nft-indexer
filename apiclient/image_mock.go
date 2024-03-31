package apiclient

import (
	"context"
)

type MockImageClient struct {
	URLToFilename map[string]string
}

func NewMockImageClient(URLToFilename map[string]string) (c MockImageClient) {
	c.URLToFilename = URLToFilename
	return
}

func (c *MockImageClient) Download(ctx context.Context, imageURL, tag string, overwrite bool) (location string, err error) {
	filename, ok := c.URLToFilename[imageURL]
	if !ok {
		return "", ErrImageNotFound
	}
	return DefaultBaseImageURL + "/" + tag + "/" + filename, nil
}
