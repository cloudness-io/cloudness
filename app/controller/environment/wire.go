package environment

import (
	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/volume"
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
	appCtrl *application.Controller,
	volumeCtrl *volume.Controller,
	environmentStore store.EnvironmentStore,
) *Controller {
	return NewController(
		tx,
		appCtrl,
		volumeCtrl,
		environmentStore,
	)
}
