package client

import (
	"context"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/types"
)

type RunnerClient interface {
	// GetPrimaryServerUID returns the primary server UID.
	GetPrimaryServerUID(ctx context.Context) (int64, error) //NOTE: temporary

	// Config returns the runner config for the server.
	Config(ctx context.Context, serverUID int64) (*types.RunnerConfig, error)

	// Request requests the next available deployment for execution.
	Request(ctx context.Context) (*types.Deployment, error)

	// Accept accepts a deployment for execution.
	Accept(ctx context.Context, deployment *types.Deployment) error

	// Init returns all the information for the deployment to be executed.
	Init(ctx context.Context, deploymentID int64) (*pipeline.RunnerContext, error)

	// Update updates the deployment with new information.
	Update(ctx context.Context, deployment *types.Deployment) error

	// Watch watches for deployment cancellation events.
	Watch(ctx context.Context, deploymentID int64) (bool, error)

	// Stream streams real time logs to backend
	Stream(ctx context.Context, deploymentID int64, logs []*types.LogLine) error

	// Upload uploads the full logs to server
	Upload(ctx context.Context, deploymentID int64, logs []*types.LogLine) error

	//Metrics
	// UploadMetrics uploads the full metrics to server
	UploadMetrics(ctx context.Context, metrics []*types.AppMetrics) error
}
