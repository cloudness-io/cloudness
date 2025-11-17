package types

type Template struct {
	ID       int64         `db:"template_id"         json:"-"`
	Slug     string        `db:"template_slug"       json:"slug"`
	Name     string        `db:"template_name"       json:"name"`
	Icon     string        `db:"template_icon"       json:"icon"`
	ReadMe   string        `db:"template_readme"     json:"readme"`
	Tags     string        `db:"template_tags"       json:"tags"`
	Spec     *TemplateSpec `db:"-"                   json:"spec"`
	SpecJson string        `db:"template_spec"       json:"-"`
	Created  int64         `db:"template_created"    json:"created"`
}
