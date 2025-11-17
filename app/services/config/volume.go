package config

import (
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) GetVolumeRestrictions(server *types.Server, tenant *types.Tenant) *types.VolumeRestriction {
	return &types.VolumeRestriction{
		MaxVolumes:    tenant.MaxVolumes,
		MinVolumeSize: helpers.Max(tenant.MinVolumeSize, server.VolumeMinSize),
		MaxVolumeSize: tenant.MaxVolumeSize,
	}
}
