package volume

import (
	"context"

	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types"
)

type Controller struct {
	configSvc   *config.Service
	volumeStore store.VolumeStore
}

func NewController(configSvc *config.Service, volumeStore store.VolumeStore) *Controller {
	return &Controller{
		configSvc:   configSvc,
		volumeStore: volumeStore,
	}
}

func (c *Controller) findByUID(ctx context.Context, tenantID, projectID, envID, volumeUID int64) (*types.Volume, error) {
	return c.volumeStore.FindByUID(ctx, tenantID, projectID, envID, volumeUID)
}
