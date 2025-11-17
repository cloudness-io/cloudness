package deployment

import (
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(deploymentStore store.DeploymentStore) *Controller {
	return NewController(deploymentStore)
}
