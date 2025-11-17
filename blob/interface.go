package blob

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrNotSupported = errors.New("not supported")
)

type Store interface {
	// Upload uploads a file to the blob store.
	Upload(ctx context.Context, file io.Reader, filePath string) error

	// GetSignedURL returns the URL for a file in the blob store.
	GetSignedURL(ctx context.Context, filePath string) (string, error)

	// Download returns a reader for a file in the blob store.
	Download(ctx context.Context, filePath string) (io.ReadCloser, error)
}
