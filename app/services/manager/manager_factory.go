package manager

import (
	"fmt"

	"github.com/cloudness-io/cloudness/app/services/manager/kube"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

type ManagerFactory struct {
	managers map[enum.ServerType]ServerManager
}

func NewManagerFactory(k8sManager *kube.K8sManager) ManagerFactory {
	managers := make(map[enum.ServerType]ServerManager)
	managers[enum.ServerTypeK8s] = k8sManager
	return ManagerFactory{
		managers: managers,
	}
}

func (f *ManagerFactory) GetServerManager(server *types.Server) (ServerManager, error) {
	val := f.managers[server.Type]
	if val == nil {
		return nil, fmt.Errorf("unknown server manager type: %s", server.Type)
	}

	return val, nil
}

func (f *ManagerFactory) GetServerManagerByType(serverType enum.ServerType) (ServerManager, error) {
	val := f.managers[serverType]
	if val == nil {
		return nil, fmt.Errorf("unknown server manager type: %s", serverType)
	}

	return val, nil
}
