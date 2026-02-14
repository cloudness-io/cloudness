package types

import (
	"github.com/cloudness-io/cloudness/types/enum"
)

type Volume struct {
	ID             int64  `db:"volume_id"              json:"-"`
	UID            int64  `db:"volume_uid"             json:"uid"`
	TenantID       int64  `db:"volume_tenant_id"       json:"-"`
	ProjectID      int64  `db:"volume_project_id"      json:"-"`
	EnvironmentID  int64  `db:"volume_environment_id"  json:"-"`
	EnvironmentUID int64  `db:"volume_environment_uid" json:"environment_uid"`
	ServerID       int64  `db:"volume_server_id"       json:"-"`
	ApplicaitonID  *int64 `db:"volume_application_id"  json:"application_id"`
	Name           string `db:"volume_name"            json:"name"`
	Slug           string `db:"volume_slug"            json:"slug"`
	ParentSlug     string `db:"volume_parent_slug"     json:"parent_slug"`
	MountPath      string `db:"volume_mount_path"      json:"mount_path"`
	HostPath       string `db:"volume_host_path"       json:"host_path"`
	Size           int64  `db:"volume_size"            json:"size"`

	Created int64  `db:"volume_created" json:"created"`
	Updated int64  `db:"volume_updated" json:"updated"`
	Deleted *int64 `db:"volume_deleted" json:"deleted"`
}

type VolumeFilter struct {
	ListQueryFilter
	TenantID          *int64          `json:"tenant_id,omitempty"`
	ProjectID         *int64          `json:"project_id,omitempty"`
	EnvironmentID     *int64          `json:"environment_id,omitempty"`
	ApplicationID     *int64          `json:"application_id,omitempty"`
	ServerID          *int64          `json:"server_id,omitempty"`
	Sort              enum.VolumeAttr `json:"sort"`
	Order             enum.Order      `json:"order"`
	DeletedAt         *int64          `json:"deleted_at,omitempty"`
	DeletedBeforeOrAt *int64          `json:"deleted_before_or_at,omitempty"`
}

// input
type VolumeCreateInput struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	MountPath string  `json:"mountPath"`
	Size      int64   `json:"size,string"`
	Server    *Server `json:"-"`
}

// helpers
func (v *Volume) GetIdentifierStr() string {
	return v.Slug
}

func (v *Volume) ToInput() *VolumeCreateInput {
	return &VolumeCreateInput{
		Name:      v.Name,
		Slug:      v.Slug,
		MountPath: v.MountPath,
		Size:      v.Size,
	}
}

func (v *Volume) FromInput(in *VolumeCreateInput) {
	v.Name = in.Name
	v.MountPath = in.MountPath
	v.Size = in.Size
}
