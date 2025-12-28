package background

import (
	"context"
	"sync"
	"time"

	"github.com/cloudness-io/cloudness/lock"

	"github.com/rs/zerolog/log"
)

type BackgroundService interface {
	Start(ctx context.Context) error
	Name() string
	WaitFor() time.Duration
}

type Service struct {
	mxManager          lock.MutexManager
	BackgroundServices []BackgroundService
}

func New(
	mxManager lock.MutexManager,
) *Service {
	return &Service{
		mxManager:          mxManager,
		BackgroundServices: []BackgroundService{},
	}
}

func (s *Service) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, service := range s.BackgroundServices {
		wg.Go(func() {
			s.startService(ctx, service)
		})
	}
	wg.Wait()

	return nil
}

func (s *Service) startService(ctx context.Context, service BackgroundService) error {
	mx, err := s.mxManager.NewMutex(service.Name())
	if err != nil {
		return err
	}

	defer func() {
		if err := mx.Unlock(ctx); err != nil {
			log.Error().Err(err).Msg("failed to release global lock after starting background service")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !mx.IsHeld(ctx) {
				if err := service.Start(ctx); err != nil {
					log.Error().Err(err).Msg("failed to start background service")
					return err
				}
			}
			time.Sleep(service.WaitFor())
		}
	}
}
