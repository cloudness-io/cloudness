package engine

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/types"
)

type (
	Engine interface {
		// Type returns the type of executor
		Type() string

		IsAvailable(config types.RunnerConfig) bool

		Load(ctx context.Context, config types.RunnerConfig) (*EngineInfo, error)

		ListIncomplete(ctx context.Context) ([]int64, error)

		Setup(ctx context.Context, pCtx *pipeline.RunnerContext) error

		StartStep(ctx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) error

		TailStep(ctx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) (io.ReadCloser, error)

		WaitStep(ctx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) (*State, error)

		Destroy(ctx context.Context, pCtx *pipeline.RunnerContext) error
	}

	State struct {
		// Container exit code
		ExitCode int
		// Container exited, true or false
		Exited bool
		// Container is OOMKilled, true or false
		OOMKilled bool
		// Container error
		Error error
	}
)

type EngineInfo struct{}

func FindEngine(engines []Engine, config types.RunnerConfig) (Engine, error) {
	if config.Engine == "" {
		for _, engine := range engines {
			if engine.IsAvailable(config) {
				return engine, nil
			}
		}
		// return engines[0], nil
		return nil, fmt.Errorf("unable to detect any backend engine")
	}

	for _, engine := range engines {
		if engine.Type() == config.Engine && engine.IsAvailable(config) {
			return engine, nil
		}
	}
	return nil, fmt.Errorf("backend engine %s is not available", config.Engine)

}
