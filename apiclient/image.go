package apiclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrImageNotFound      = errors.New("err image not found")
	ErrImageAlreadyExists = errors.New("err image already exists")
)

const (
	LocalImageDownloadClientTimeout = int64(20) // 20s
	DefaultBaseDownloadLocation     = "./static"
	DefaultBaseImageURL             = "http://127.0.0.1:3000/image"
	CommonTag                       = "common"
)

type TokenURI struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

var _ ImageDownloadClient = (*LocalImageDownloadClient)(nil)

type LocalImageDownloadClient struct {
	c http.Client

	baseDownloadLocation string
	baseImageURL         string
}

func NewLocalImageDownloadClient(ctx context.Context, baseDownloadLocation, baseImageURL string) (c LocalImageDownloadClient, err error) {
	c.c = http.Client{
		Timeout: time.Duration(LocalImageDownloadClientTimeout * int64(time.Second)),
	}

	if baseDownloadLocation != "" {
		c.baseDownloadLocation = baseDownloadLocation
	} else {
		c.baseDownloadLocation = DefaultBaseDownloadLocation
	}

	if baseImageURL != "" {
		c.baseImageURL = baseImageURL
	} else {
		c.baseImageURL = DefaultBaseImageURL
	}

	return
}

func (c *LocalImageDownloadClient) Download(ctx context.Context, imageURL, tag string, overwrite bool) (location string, err error) {
	var (
		localPath = filepath.Join(c.baseDownloadLocation, tag)
		filename  = filepath.Base(imageURL)
		fullPath  = filepath.Join(localPath, filename)
		req       *http.Request
		resp      *http.Response
	)

	if !overwrite && fileExists(fullPath) {
		err = fmt.Errorf("%w: %s", ErrImageAlreadyExists, fullPath)
		return
	}

	if req, err = http.NewRequestWithContext(ctx, "GET", imageURL, nil); err != nil {
		err = fmt.Errorf("at http.NewRequestWithContext: %w", err)
		return
	}

	if resp, err = c.c.Do(req); err != nil {
		if resp.StatusCode == http.StatusNotFound {
			// 404 error
			err = fmt.Errorf("%w: %w", ErrImageNotFound, err)
		}
		err = fmt.Errorf("failed to downlad image from. %s: %w", imageURL, err)
		return
	}
	defer resp.Body.Close()

	// Create the directory if it doesn't exist
	if err = os.MkdirAll(localPath, os.ModePerm); err != nil {
		panic(err)
	}

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		err = fmt.Errorf("failed to create image file. %s: %w", fullPath, err)
		return
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		err = fmt.Errorf("failed to copy image to file. %s: %w", fullPath, err)
	}

	location = fmt.Sprintf("%s/%s/%s", c.baseImageURL, tag, filename)
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
