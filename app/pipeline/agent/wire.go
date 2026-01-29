package agent

import (
	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/app/services/manager"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideAgent,
)

func ProvideAgent(client client.RunnerClient, factory manager.ManagerFactory) *Agent {
	return New(client, factory)
}
