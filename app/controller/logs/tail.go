package logs

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Tail(ctx context.Context, deploymentID int64) (<-chan *types.LogLine, <-chan error, error) {
	linec, errc := c.logStream.Tail(ctx, deploymentID)
	return linec, errc, nil
}
