package types

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
)

func (t *TemplateSpec) ToTemplate() (*Template, error) {
	tmpl := &Template{
		Slug:    helpers.Normalize(t.Name),
		Name:    t.Name,
		Icon:    t.Icon,
		ReadMe:  t.Readme,
		Tags:    strings.Join(t.Tags, ","),
		Spec:    t,
		Created: time.Now().UTC().UnixMilli(),
	}

	specJson, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	tmpl.SpecJson = string(specJson)

	return tmpl, nil
}
