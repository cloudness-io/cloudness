package background

import (
	"github.com/cloudness-io/cloudness/lock"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
)

func ProvideService(
	mxManager lock.MutexManager,
) *Service {
	return New(
		mxManager,
	)
}
