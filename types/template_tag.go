package types

type Tag struct {
	ID      int64  `db:"tag_id"        json:"-"`
	Slug    string `db:"tag_slug"      json:"slug"`
	Name    string `db:"tag_name"      json:"name"`
	Created int64  `db:"tag_created"   json:"created"`
}

type TemplateTag struct {
	TemplateID int64 `db:"template_id" json:"-"`
	TagID      int64 `db:"tag_id"      json:"-"`
}
