package config

import "github.com/cloudness-io/cloudness/types"

func (s *Service) GetAppRestrictions(tenant *types.Tenant) *types.ApplicationRestrction {
	return &types.ApplicationRestrction{
		MaxInstance: tenant.MaxInstances,
		MaxCPU:      tenant.MaxCPUPerApp,
		MaxMemory:   tenant.MaxMemoryPerApp,
	}
}
