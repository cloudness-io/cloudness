package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type queue struct {
	sync.Mutex
	globMx lock.Mutex

	ready    chan struct{}
	paused   bool
	interval time.Duration
	store    store.DeploymentStore
	managers map[*manager]struct{}
	ctx      context.Context
}

// newQueue returns a new Queue backed by the build datastore.
func newQueue(store store.DeploymentStore, lock lock.MutexManager) (*queue, error) {
	const lockKey = "build_queue"
	mx, err := lock.NewMutex(lockKey)
	if err != nil {
		return nil, err
	}
	q := &queue{
		store:    store,
		globMx:   mx,
		ready:    make(chan struct{}, 1),
		managers: map[*manager]struct{}{},
		interval: 15 * time.Second,
		ctx:      context.Background(),
	}
	go func() {
		if err := q.start(); err != nil {
			log.Err(err).Msg("queue start failed")
		}
	}()

	return q, nil
}

func (q *queue) Schedule(_ context.Context, _ *types.Deployment) error {
	select {
	case q.ready <- struct{}{}:
	default:
	}
	return nil
}

func (q *queue) Pause(_ context.Context) error {
	q.Lock()
	q.paused = true
	q.Unlock()
	return nil
}

func (q *queue) Request(ctx context.Context, params Filter) (*types.Deployment, error) {
	w := &manager{
		channel: make(chan *types.Deployment),
		done:    ctx.Done(),
	}
	q.Lock()
	q.managers[w] = struct{}{}
	q.Unlock()

	select {
	case q.ready <- struct{}{}:
	default:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b := <-w.channel:
		return b, nil
	}
}

//nolint:gocognit // refactor if needed.
func (q *queue) signal(ctx context.Context) error {
	if err := q.globMx.Lock(ctx); err != nil {
		return err
	}
	defer func() {
		if err := q.globMx.Unlock(ctx); err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to release global lock after signaling")
		}
	}()

	q.Lock()
	count := len(q.managers)
	pause := q.paused
	q.Unlock()
	if pause {
		return nil
	}
	if count == 0 {
		return nil
	}
	items, err := q.store.ListIncomplete(ctx)
	log.Ctx(ctx).Error().Msg("scheduler: list incomplete deployments")
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("scheduler: failed to get incomplete deployments")
		return err
	}

	q.Lock()
	defer q.Unlock()
	for _, item := range items {
		if item.Status == enum.DeploymentStatusRunning {
			continue
		}

	loop:
		for m := range q.managers {
			select {
			case m.channel <- item:
			case <-m.done:
			case <-time.After(q.interval):
			}

			delete(q.managers, m)
			break loop
		}
	}
	return nil
}

func (q *queue) start() error {
	for {
		select {
		case <-q.ctx.Done():
			return q.ctx.Err()
		case <-q.ready:
			if err := q.signal(q.ctx); err != nil {
				// don't return, only log error
				log.Ctx(q.ctx).Err(err).Msg("scheduler: failed to signal on ready")
			}
		case <-time.After(q.interval):
			if err := q.signal(q.ctx); err != nil {
				// don't return, only log error
				log.Ctx(q.ctx).Err(err).Msg("scheduler: failed to signal on interval")
			}
		}
	}
}

type manager struct {
	labels  map[string]string
	channel chan *types.Deployment
	done    <-chan struct{}
}
