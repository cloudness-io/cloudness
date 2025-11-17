package config

import (
	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) GetTenantDefaults() types.TenantConfig {
	return s.config.TenantConfig
}

func (s *Service) GetTenantRestrictions(tenant *types.Tenant) *types.TenantRestrictions {
	return &types.TenantRestrictions{
		AllowAdminToModify: tenant.AllowAdminToModify,
		MaxProjects:        tenant.MaxProjects,
		MaxApps:            tenant.MaxApps,
		MaxInstances:       tenant.MaxInstances,
		MaxCPU:             tenant.MaxCPUPerApp,
		MaxMemory:          tenant.MaxMemoryPerApp,
		MaxVolumes:         tenant.MaxVolumes,
		MinVolumeSize:      tenant.MinVolumeSize,
		MaxVolumeSize:      tenant.MaxVolumeSize,
	}
}
