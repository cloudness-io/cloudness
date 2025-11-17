package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/pipeline/canceler"
	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/services/schema"
	"github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Controller struct {
	tx               dbtx.Transactor
	configSvc        *config.Service
	schemaSvc        *schema.Service
	specSvc          *spec.Service
	applicationStore store.ApplicationStore
	serverCtrl       *server.Controller
	varCtrl          *variable.Controller
	gitPublicCtrl    *gitpublic.Controller
	volumeCtrl       *volume.Controller
	triggerer        triggerer.Triggerer
	canceler         canceler.Canceler
	manager          manager.ManagerFactory
}

func NewController(
	tx dbtx.Transactor,
	configSvc *config.Service,
	schemaSvc *schema.Service,
	specSvc *spec.Service,
	applicationStore store.ApplicationStore,
	serverCtrl *server.Controller,
	varCtrl *variable.Controller,
	gitPublicCtrl *gitpublic.Controller,
	volumeCtrl *volume.Controller,
	triggerer triggerer.Triggerer,
	canceler canceler.Canceler,
	manager manager.ManagerFactory,
) *Controller {
	return &Controller{
		tx:               tx,
		configSvc:        configSvc,
		schemaSvc:        schemaSvc,
		specSvc:          specSvc,
		applicationStore: applicationStore,
		serverCtrl:       serverCtrl,
		varCtrl:          varCtrl,
		gitPublicCtrl:    gitPublicCtrl,
		volumeCtrl:       volumeCtrl,
		triggerer:        triggerer,
		canceler:         canceler,
		manager:          manager,
	}
}

func (c *Controller) findByUID(ctx context.Context, tenantID, projectID, environmentID, applicationUID int64) (*types.Application, error) {
	return c.applicationStore.FindByUID(ctx, tenantID, projectID, environmentID, applicationUID)
}

func (c *Controller) updateWithTx(ctx context.Context, dto *createOrUpdateDto) (*types.Application, error) {
	application := dto.Application
	var err error

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		application, err = c.updateWithoutTx(ctx, dto)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return application, err
}

func (c *Controller) updateWithoutTx(ctx context.Context, dto *createOrUpdateDto) (*types.Application, error) {
	application := dto.Application
	err := c.schemaSvc.ValidateApplication(ctx, application.Spec)
	if err != nil {
		return nil, err
	}

	//Validate restrictions
	err = c.validateRestrictions(dto)
	if err != nil {
		return nil, err
	}

	application.DeploymentStatus = enum.ApplicationDeploymentStatusNeedsDeployment
	//update application spec json
	err = application.UpdateSpecJSON()
	if err != nil {
		return nil, err
	}

	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return nil, err
	}

	application, err = c.applicationStore.UpdateSpec(ctx, application)
	if err != nil {
		return nil, err
	}
	if err := c.varCtrl.UpdateDefaultVariables(ctx, server, dto.Tenant, dto.Project, dto.Environment, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (c *Controller) validateRestrictions(dto *createOrUpdateDto) error {
	deploySpec := dto.Application.Spec.Deploy
	restrictions := c.configSvc.GetTenantRestrictions(dto.Tenant)

	err := check.NewValidationErrors()
	if deploySpec.CPU > restrictions.MaxCPU {
		err.AddValidationError("cpu", check.NewValidationErrorf("CPU is above max allowed limit"))
	}
	if deploySpec.MaxReplicas > restrictions.MaxInstances {
		err.AddValidationError("maxReplicas", check.NewValidationErrorf("Max Replicas is above max allowed limit"))
	}
	if deploySpec.Memory > restrictions.MaxMemory {
		err.AddValidationError("memory", check.NewValidationErrorf("Memory is above max allowed limit"))
	}

	if err.HasError() {
		return err
	}
	return nil
}
