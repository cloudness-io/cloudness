package logstream

import (
	"fmt"
)

func getStreamId(deploymentId int64) string {
	return fmt.Sprintf("%d-%s", deploymentId)
}
