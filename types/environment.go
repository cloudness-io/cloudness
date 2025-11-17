package types

import "github.com/cloudness-io/cloudness/types/enum"

type Environment struct {
	ID        int64  `db:"environment_id"           json:"-"`
	TenantID  int64  `db:"environment_tenant_id"    json:"-"`
	ProjectID int64  `db:"environment_project_id"   json:"-"`
	UID       int64  `db:"environment_uid"          json:"uid"`
	Name      string `db:"environment_name"         json:"name"`
	CreateBy  int64  `db:"environment_created_by"   json:"-"`

	Created int64  `db:"environment_created"                json:"created"`
	Updated int64  `db:"environment_updated"                json:"updated"`
	Deleted *int64 `db:"environment_deleted"                json:"deleted"`
}

// EnvironmentFilter stores projects query parameters.
type EnvironmentFilter struct {
	ListQueryFilter
	TenantID          *int64               `json:"tenant_id,omitempty"`
	ProjectID         *int64               `json:"project_id,omitempty"`
	Sort              enum.EnvironmentAttr `json:"sort"`
	Order             enum.Order           `json:"order"`
	DeletedAt         *int64               `json:"deleted_at,omitempty"`
	DeletedBeforeOrAt *int64               `json:"deleted_before_or_at,omitempty"`
}
