package client

import (
	"github.com/cloudness-io/cloudness/app/pipeline/manager"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideRunnerClient,
)

func ProvideRunnerClient(manager manager.RunnerManager) RunnerClient {
	return NewEmbeddedClient(manager)
}
