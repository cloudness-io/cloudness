package template

import (
	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/app/services/schema"
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(templateStore store.TemplateStore, appCtrl *application.Controller, schemaSvc *schema.Service, triggerer triggerer.Triggerer) *Controller {
	return NewController(templateStore, appCtrl, schemaSvc, triggerer)
}
