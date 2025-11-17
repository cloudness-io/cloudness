package environment

import (
	"context"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type Controller struct {
	tx               dbtx.Transactor
	appCtrl          *application.Controller
	volumeCtrl       *volume.Controller
	environmentStore store.EnvironmentStore
}

func NewController(
	tx dbtx.Transactor,
	appCtrl *application.Controller,
	volumeCtrl *volume.Controller,
	environmentStore store.EnvironmentStore,
) *Controller {
	return &Controller{
		tx:               tx,
		appCtrl:          appCtrl,
		volumeCtrl:       volumeCtrl,
		environmentStore: environmentStore,
	}
}

func (c *Controller) FindByUID(ctx context.Context, projectID, environmentUID int64) (*types.Environment, error) {
	return c.environmentStore.FindByUID(ctx, projectID, environmentUID)
}

func (c *Controller) findEnvironmentByID(ctx context.Context, environmentID int64) (*types.Environment, error) {
	return c.environmentStore.Find(ctx, environmentID)
}

func (c *Controller) sanitizeCreateInput(in *CreateEnvironmentInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
