package job

import (
	"context"

	"github.com/cloudness-io/cloudness/lock"
)

func globalLock(ctx context.Context, manager lock.MutexManager) (lock.Mutex, error) {
	const lockKey = "jobs"
	mx, err := manager.NewMutex(lockKey)
	if err != nil {
		return nil, err
	}

	err = mx.Lock(ctx)

	return mx, err
}
