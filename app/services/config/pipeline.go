package config

import "context"

func (s *Service) Pipeline(ctx context.Context) (*PipelineConfig, error) {
	instance, err := s.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	c := &PipelineConfig{}

	c.PushRegistryURL = "registry.cloudness.svc.cluster.local:5000"
	c.PullRegistryURL = "localhost:30050"
	if instance.RegistryMirrorEnabled {
		c.MirrorRegistryEnabled = true
		c.MirrorRegistryURL = "registry-mirror.cloudness.svc.cluster.local:5000"
	}
	return c, nil
}
