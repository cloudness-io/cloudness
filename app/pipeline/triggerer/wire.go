package triggerer

import (
	"github.com/cloudness-io/cloudness/app/pipeline/canceler"
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideTriggerer,
)

// ProvideTriggerer provides a triggerer which can execute builds.
func ProvideTriggerer(
	tx dbtx.Transactor,
	appStore store.ApplicationStore,
	deploymentStore store.DeploymentStore,
	scheduler scheduler.Scheduler,
	canceler canceler.Canceler,
) Triggerer {
	return New(tx, appStore, deploymentStore, scheduler, canceler)
}
