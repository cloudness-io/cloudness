package types

import "github.com/cloudness-io/cloudness/types/enum"

type Tenant struct {
	ID                 int64   `db:"tenant_id"                     json:"-"`
	UID                int64   `db:"tenant_uid"                    json:"uid"`
	Name               string  `db:"tenant_name"                   json:"name"`
	Description        string  `db:"tenant_description"            json:"description"`
	AllowAdminToModify bool    `db:"tenant_allow_admin_to_modify"  json:"allow_admin_to_modify"`
	MaxProjects        int64   `db:"tenant_max_projects"           json:"max_projects"`
	MaxApps            int64   `db:"tenant_max_apps"               json:"max_apps"`
	MaxInstances       int64   `db:"tenant_max_instances"          json:"max_instances"`
	MaxCPUPerApp       int64   `db:"tenant_max_cpu_per_app"        json:"max_cpu_per_app"`
	MaxMemoryPerApp    float64 `db:"tenant_max_memory_per_app"     json:"max_memory_per_app"`
	MaxVolumes         int64   `db:"tenant_max_volumes"            json:"max_volumes"`
	MinVolumeSize      int64   `db:"tenant_min_volume_size"        json:"min_volume_size"`
	MaxVolumeSize      int64   `db:"tenant_max_volume_size"        json:"max_volume_size"`
	CreateBy           int64   `db:"tenant_created_by"             json:"-"`

	Created int64  `db:"tenant_created"                json:"created"`
	Updated int64  `db:"tenant_updated"                json:"updated"`
	Deleted *int64 `db:"tenant_deleted"                json:"deleted"`
}

// TenantFilter stores projects query parameters.
type TenantFilter struct {
	ListQueryFilter
	Sort              enum.TenantAttr `json:"sort"`
	Order             enum.Order      `json:"order"`
	DeletedAt         *int64          `json:"deleted_at,omitempty"`
	DeletedBeforeOrAt *int64          `json:"deleted_before_or_at,omitempty"`
}
