package types

import "github.com/cloudness-io/cloudness/types/enum"

type Project struct {
	ID          int64  `db:"project_id"           json:"-"`
	TenantID    int64  `db:"project_tenant_id"    json:"-"`
	UID         int64  `db:"project_uid"          json:"uid"`
	Name        string `db:"project_name"         json:"name"`
	Slug        string `db:"project_slug"         json:"slug"`
	Description string `db:"project_description"  json:"description"`
	CreateBy    int64  `db:"project_created_by"   json:"-"`

	Created int64  `db:"project_created"     json:"created"`
	Updated int64  `db:"project_updated"     json:"updated"`
	Deleted *int64 `db:"project_deleted"     json:"deleted"`
}

// ProjectFilter stores projects query parameters.
type ProjectFilter struct {
	ListQueryFilter
	TenantID          *int64           `json:"tenant_id,omitempty"`
	PrincipalID       *int64           `json:"principal_id,omitempty"`
	Sort              enum.ProjectAttr `json:"sort"`
	Order             enum.Order       `json:"order"`
	DeletedAt         *int64           `json:"deleted_at,omitempty"`
	DeletedBeforeOrAt *int64           `json:"deleted_before_or_at,omitempty"`
}
