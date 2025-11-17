package application

import (
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

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
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
	return NewController(
		tx,
		configSvc,
		schemaSvc,
		specSvc,
		applicationStore,
		serverCtrl,
		varCtrl,
		gitPublicCtrl,
		volumeCtrl,
		triggerer,
		canceler,
		manager,
	)
}
