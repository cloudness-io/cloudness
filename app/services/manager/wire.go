package manager

import (
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/manager/kube"
	"github.com/cloudness-io/cloudness/app/services/manager/vm"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideManagerFactory,
	ProvideK8sManager,
	ProvideVmManager,
)

func ProvideManagerFactory(k8sManager *kube.K8sManager, vmManager *vm.VmManager) ManagerFactory {
	return NewManagerFactory(k8sManager, vmManager)
}

func ProvideK8sManager(configSvc *config.Service) *kube.K8sManager {
	return kube.NewK8sManager(configSvc)
}

func ProvideVmManager(configSvc *config.Service) *vm.VmManager {
	return vm.NewVmManager(configSvc)
}
