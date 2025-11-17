package scheduler

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/types"
)

// Filter provides filter criteria to limit stages requested
// from the scheduler.
type Filter struct {
	Labels map[string]string
}

// Scheduler schedules deployments
type Scheduler interface {
	// Schedule schedules a deployment for execution
	Schedule(ctx context.Context, deployment *types.Deployment) error

	// Request requests the next stage scheduled for execution.
	Request(ctx context.Context, filter Filter) (*types.Deployment, error)

	// Cancel cancels scheduled or running jobs associated
	// with the parent build ID.
	Cancel(context.Context, int64) error

	// Cancelled blocks and listens for a cancellation event and
	// returns true if the build has been cancelled.
	Cancelled(context.Context, int64) (bool, error)
}

type scheduler struct {
	*queue
	*canceler
}

// newScheduler provides an instance of a scheduler with cancel abilities.
func newScheduler(deploymentStore store.DeploymentStore, lock lock.MutexManager) (Scheduler, error) {
	q, err := newQueue(deploymentStore, lock)
	if err != nil {
		return nil, err
	}
	return scheduler{
		q,
		newCanceler(),
	}, nil
}
