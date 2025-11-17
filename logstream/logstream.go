package logstream

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

type LogStream interface {
	// Create creates a log stream for the deployment id
	Create(context.Context, int64) error

	// Delete deletes a log stream for deployment id
	Delete(context.Context, int64) error

	// Write writes to the log stream
	Write(context.Context, int64, *types.LogLine) error

	// Tail tails the log stream
	Tail(context.Context, int64) (<-chan *types.LogLine, <-chan error)
}
