package project

import (
	"context"

	"github.com/cloudness-io/cloudness/app/sse"
)

func (c *Controller) Events(ctx context.Context, projectID int64) (<-chan *sse.Event, <-chan error, func(context.Context) error, error) {

	chEvents, chError, sseCancel := c.sseStremer.Stream(ctx, projectID)

	return chEvents, chError, sseCancel, nil
}
