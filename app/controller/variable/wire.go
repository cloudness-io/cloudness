package variable

import (
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(variableStore store.VariableStore) *Controller {
	return NewController(variableStore)
}
