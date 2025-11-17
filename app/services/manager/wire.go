package manager

import (
	"github.com/cloudness-io/cloudness/app/services/manager/kube"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideManagerFactory,
	ProvideK8sManager,
)

func ProvideK8sManager() *kube.K8sManager {
	return kube.NewK8sManager()
}

func ProvideManagerFactory(k8sManager *kube.K8sManager) ManagerFactory {
	return NewManagerFactory(k8sManager)
}
