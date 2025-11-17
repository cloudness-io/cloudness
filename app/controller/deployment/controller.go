package deployment

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/store"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

type Controller struct {
	deploymentStore store.DeploymentStore
}

func NewController(deploymentStore store.DeploymentStore) *Controller {
	return &Controller{
		deploymentStore: deploymentStore,
	}
}

func (c *Controller) findByUID(ctx context.Context, applicationID, deploymentUID int64) (*types.Deployment, error) {
	deployment, err := c.deploymentStore.FindByUID(ctx, applicationID, deploymentUID)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	return deployment, nil
}
