package vm

import "github.com/cloudness-io/cloudness/app/services/config"

type VmManager struct {
	configSvc *config.Service
}

func NewVmManager(configSvc *config.Service) *VmManager {
	return &VmManager{configSvc: configSvc}
}
