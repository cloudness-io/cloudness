package project

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type Controller struct {
	tx                     dbtx.Transactor
	configSvc              *config.Service
	userCtrl               *user.Controller
	envCtrl                *environment.Controller
	projectStore           store.ProjectStore
	projectMembershipStore store.ProjectMembershipStore
	tenantMembershipStore  store.TenantMembershipStore
	sseStremer             sse.Streamer
}

func NewController(tx dbtx.Transactor,
	configSvc *config.Service,
	userCtrl *user.Controller,
	envCtrl *environment.Controller,
	projectStore store.ProjectStore,
	projectMembershipStore store.ProjectMembershipStore,
	tenantMembershipStore store.TenantMembershipStore,
	sseStremer sse.Streamer,
) *Controller {
	return &Controller{
		tx:                     tx,
		configSvc:              configSvc,
		userCtrl:               userCtrl,
		envCtrl:                envCtrl,
		projectStore:           projectStore,
		projectMembershipStore: projectMembershipStore,
		tenantMembershipStore:  tenantMembershipStore,
		sseStremer:             sseStremer,
	}
}

func (c *Controller) findByID(ctx context.Context, projectID int64) (*types.Project, error) {
	project, err := c.projectStore.Find(ctx, projectID)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	return project, nil
}

func (c *Controller) findByUID(ctx context.Context, tenantID int64, projectUID int64) (*types.Project, error) {
	project, err := c.projectStore.FindByUID(ctx, tenantID, projectUID)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	return project, nil
}

func (c *Controller) sanitizeCreateInput(in *CreateProjectInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if err := check.Description(in.Description); err != nil {
		errors.AddValidationError("description", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
