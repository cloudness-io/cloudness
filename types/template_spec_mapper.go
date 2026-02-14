package types

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
)

func (t *TemplateSpec) ToTemplate() (*Template, error) {
	normalizedTags := normalizeTags(t.Tags)
	t.Tags = normalizedTags

	tmpl := &Template{
		Slug:    helpers.Normalize(t.Name),
		Name:    t.Name,
		Icon:    t.Icon,
		ReadMe:  t.Readme,
		Tags:    normalizedTags,
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

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(tags))

	for _, tag := range tags {
		clean := strings.TrimSpace(tag)
		if clean == "" {
			continue
		}

		key := strings.ToLower(clean)
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		normalized = append(normalized, clean)
	}

	return normalized
}
