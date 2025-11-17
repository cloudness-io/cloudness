package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Find(ctx context.Context, deploymentID int64) ([]*types.LogLine, error) {
	rc, err := c.logStore.Find(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("could not find logs: %w", err)
	}
	defer rc.Close()

	lines := []*types.LogLine{}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read from buffer: %w", err)
	}

	err = json.Unmarshal(buf.Bytes(), &lines)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}

	return lines, nil
}
