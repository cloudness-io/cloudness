package manager

import (
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/manager/kube"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideManagerFactory,
	ProvideK8sManager,
)

func ProvideK8sManager(configSvc *config.Service) *kube.K8sManager {
	return kube.NewK8sManager(configSvc)
}

func ProvideManagerFactory(k8sManager *kube.K8sManager) ManagerFactory {
	return NewManagerFactory(k8sManager)
}
