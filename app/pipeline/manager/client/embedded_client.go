package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/pipeline/manager"
	"github.com/cloudness-io/cloudness/types"
)

type embeddedClient struct {
	manager manager.RunnerManager
}

var _ RunnerClient = (*embeddedClient)(nil)

func NewEmbeddedClient(manager manager.RunnerManager) *embeddedClient {
	return &embeddedClient{
		manager: manager,
	}
}

func (e *embeddedClient) GetPrimaryServerUID(ctx context.Context) (int64, error) {
	return e.manager.GetPrimaryServerUID(ctx)
}

func (e *embeddedClient) Config(ctx context.Context, serverUID int64) (*types.RunnerConfig, error) {
	return e.manager.Config(ctx, serverUID)
}

func (e *embeddedClient) Request(ctx context.Context) (*types.Deployment, error) {
	return e.manager.Request(ctx)
}

func (e *embeddedClient) Accept(ctx context.Context, deployment *types.Deployment) error {
	d, err := e.manager.Accept(ctx, deployment.ID, deployment.Machine)
	if err != nil {
		return err
	}
	*deployment = *d
	return nil
}

func (e *embeddedClient) Init(ctx context.Context, deploymentID int64) (*pipeline.RunnerContext, error) {
	return e.manager.Details(ctx, deploymentID)
}

func (e *embeddedClient) Update(ctx context.Context, deployment *types.Deployment) error {
	d, err := e.manager.Update(ctx, deployment)
	if err != nil {
		return err
	}
	*deployment = *d
	return nil
}

func (e *embeddedClient) Watch(ctx context.Context, deploymentID int64) (bool, error) {
	return e.manager.Watch(ctx, deploymentID)
}

func (e *embeddedClient) Stream(ctx context.Context, deploymentID int64, logs []*types.LogLine) error {
	return e.manager.Stream(ctx, deploymentID, logs)
}

func (e *embeddedClient) Upload(ctx context.Context, deploymentID int64, logs []*types.LogLine) error {
	var buffer bytes.Buffer
	out, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	_, err = buffer.Write(out)
	if err != nil {
		return err
	}
	return e.manager.Upload(ctx, deploymentID, &buffer)
}

func (e *embeddedClient) UploadMetrics(ctx context.Context, metrics []*types.AppMetrics) error {
	return e.manager.UploadMetrics(ctx, metrics)
}

func (e *embeddedClient) UploadAppStatus(ctx context.Context, status []*types.AppStatus) error {
	return e.manager.UploadAppStatus(ctx, status)
}
