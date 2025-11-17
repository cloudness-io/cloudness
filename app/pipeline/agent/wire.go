package agent

import (
	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideAgent,
)

func ProvideAgent(client client.RunnerClient) *Agent {
	return New(client)
}
