package config

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types"
)

type Service struct {
	config        *types.Config
	instanceStore store.InstanceStore
	isLocal       bool
}

func NewService(config *types.Config, instanceStore store.InstanceStore) *Service {
	s := &Service{
		config:        config,
		instanceStore: instanceStore,
	}
	if config.Environment == "local" {
		s.isLocal = true
	}
	return s
}
