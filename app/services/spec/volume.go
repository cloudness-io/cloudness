package spec

import "github.com/cloudness-io/cloudness/types"

func (s *Service) GetVolumeMounts(spec *types.ApplicationSpec) []*types.VolumeMounts {
	return spec.Volumes
}
