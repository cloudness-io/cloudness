package job

import "time"

const (
	instanceID    = "cloudness:job"
	maxRunning    = 3
	retentionTime = time.Hour * 5
)
