package template

import (
	"context"
	"encoding/json"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/app/services/schema"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/templates"
	"github.com/cloudness-io/cloudness/types"
)

type Controller struct {
	templateStore store.TemplateStore
	appCtrl       *application.Controller
	schemaSvc     *schema.Service
	triggerer     triggerer.Triggerer
}

func NewController(templateStore store.TemplateStore, appCtrl *application.Controller, schemaSvc *schema.Service, triggerer triggerer.Triggerer) *Controller {
	return &Controller{
		templateStore: templateStore,
		appCtrl:       appCtrl,
		schemaSvc:     schemaSvc,
		triggerer:     triggerer,
	}
}

func (c *Controller) Init(ctx context.Context) error {
	tmpls, err := c.readTemplatesFromFS(ctx)
	if err != nil {
		return err
	}

	return c.templateStore.UpsertMany(ctx, tmpls)
}

func (c *Controller) readTemplatesFromFS(ctx context.Context) ([]*types.Template, error) {
	ts, err := templates.List()
	if err != nil {
		return nil, err
	}

	out := make([]*types.Template, 0)
	for _, t := range ts {
		if err := c.schemaSvc.ValidateTemplate(ctx, t); err != nil {
			return nil, err
		}

		dst := new(types.TemplateSpec)
		if err := json.Unmarshal(t, dst); err != nil {
			return nil, err
		}
		tmpl, err := dst.ToTemplate()
		if err != nil {
			return nil, err
		}
		out = append(out, tmpl)
	}

	return out, nil
}
