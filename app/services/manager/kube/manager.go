package kube

import "github.com/cloudness-io/cloudness/app/services/config"

type K8sManager struct {
	configSvc *config.Service
}

func NewK8sManager(configSvc *config.Service) *K8sManager {
	return &K8sManager{configSvc: configSvc}
}
