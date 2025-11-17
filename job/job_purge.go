package job

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/lock"

	"github.com/rs/zerolog/log"
)

const (
	jobUIDPurge  = "cloudness:jobs:purge"
	jobTypePurge = "cloudness:jobs:purge"
	jobCronPurge = "15 */4 * * *" // every 4 hours at 15 minutes
)

type jobPurge struct {
	store     Store
	mxManager lock.MutexManager
	minOldAge time.Duration
}

func newJobPurge(store Store, mxManager lock.MutexManager, minOldAge time.Duration) *jobPurge {
	if minOldAge < 0 {
		minOldAge = 0
	}

	return &jobPurge{
		store:     store,
		mxManager: mxManager,
		minOldAge: minOldAge,
	}
}

func (j *jobPurge) Handle(ctx context.Context, _ string, _ ProgressReporter) (string, error) {
	mx, err := globalLock(ctx, j.mxManager)
	if err != nil {
		return "", fmt.Errorf("failed to obtain the lock to clean up old jobs")
	}

	defer func() {
		if err := mx.Unlock(ctx); err != nil {
			log.Err(err).Msg("failed to release global lock after cleaning up old jobs")
		}
	}()

	olderThan := time.Now().Add(-j.minOldAge)

	n, err := j.store.DeleteOld(ctx, olderThan)
	if err != nil {
		return "", fmt.Errorf("failed to purge old jobs")
	}

	result := "no old jobs found"
	if n > 0 {
		result = fmt.Sprintf("deleted %d old jobs", n)
	}

	return result, nil
}
