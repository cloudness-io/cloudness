package tenant

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/store"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type Controller struct {
	tx                    dbtx.Transactor
	configSvc             *config.Service
	tenantStore           store.TenantStore
	tenantMembershipStore store.TenantMembershipStore
	userCtrl              *user.Controller
	projectCtrl           *project.Controller
}

func NewController(tx dbtx.Transactor,
	configSvc *config.Service,
	tenantStore store.TenantStore,
	tenantMembershipStore store.TenantMembershipStore,
	userCtrl *user.Controller,
	projectCtrl *project.Controller,
) *Controller {
	return &Controller{
		tx:                    tx,
		configSvc:             configSvc,
		tenantStore:           tenantStore,
		tenantMembershipStore: tenantMembershipStore,
		userCtrl:              userCtrl,
		projectCtrl:           projectCtrl,
	}
}

func (c *Controller) findByID(ctx context.Context, tenantID int64) (*types.Tenant, error) {
	tenant, err := c.tenantStore.Find(ctx, tenantID)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	return tenant, nil
}

func (c *Controller) findByUID(ctx context.Context, tenantUID int64) (*types.Tenant, error) {
	tenant, err := c.tenantStore.FindByUID(ctx, tenantUID)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	return tenant, nil
}

func (c *Controller) sanitizeCreateInput(in *CreateTenantInput) error {
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
