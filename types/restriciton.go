package types

type InstanceRestrictions struct {
	AllowNewTenant bool
}

type ServerRestrictions struct {
	UnmountBeforeResize bool `json:"unmount_before_resize,string"`
}

type TenantRestrictions struct {
	AllowAdminToModify bool    `json:"allow_admin,string"`
	MaxProjects        int64   `json:"max_projects,string"`
	MaxApps            int64   `json:"max_apps,string"`
	MaxInstances       int64   `json:"max_instances,string"`
	MaxCPU             int64   `json:"max_cpu,string"`
	MaxMemory          float64 `json:"max_memory,string"`
	MaxVolumes         int64   `json:"max_volumes,string"`
	MinVolumeSize      int64   `json:"min_volume_size,string"`
	MaxVolumeSize      int64   `json:"max_volume_size,string"`
}

type ApplicationRestrction struct {
	MaxInstance int64
	MaxCPU      int64
	MaxMemory   float64
}

type VolumeRestriction struct {
	MaxVolumes    int64
	MinVolumeSize int64
	MaxVolumeSize int64
}
