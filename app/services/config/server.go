package config

import "github.com/cloudness-io/cloudness/types"

func (s *Service) GetKubeServerDefaults() types.KubeServerConfig {
	return s.config.KubeServerConfig
}

func (s *Service) GetServerRestrictions(server *types.Server) *types.ServerRestrictions {
	return &types.ServerRestrictions{
		UnmountBeforeResize: server.VolumeSupportsOnlineExpansion,
	}
}
