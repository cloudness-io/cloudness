package logs

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/logstream"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(logStore store.LogStore, logStream logstream.LogStream) *Controller {
	return NewController(logStore, logStream)
}
