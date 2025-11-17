//go:build wireinject
// +build wireinject

package main

import (
	"context"

	// "github.com/cloudness-io/cloudness/app/api/openapi"
	"github.com/cloudness-io/cloudness/app/auth/authn"
	"github.com/cloudness-io/cloudness/app/bootstrap"
	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/logs"
	"github.com/cloudness-io/cloudness/app/controller/project"
	serverCtrl "github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/pipeline/agent"
	"github.com/cloudness-io/cloudness/app/pipeline/canceler"
	"github.com/cloudness-io/cloudness/app/pipeline/manager"
	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/app/router"
	"github.com/cloudness-io/cloudness/app/server"
	"github.com/cloudness-io/cloudness/app/services"
	"github.com/cloudness-io/cloudness/app/services/cleanup"
	configSvc "github.com/cloudness-io/cloudness/app/services/config"
	dnsSvc "github.com/cloudness-io/cloudness/app/services/dns"
	githubAppSvc "github.com/cloudness-io/cloudness/app/services/githubapp"
	gitpublicSvc "github.com/cloudness-io/cloudness/app/services/gitpublic"
	managerSvc "github.com/cloudness-io/cloudness/app/services/manager"
	proxySvc "github.com/cloudness-io/cloudness/app/services/proxy"
	"github.com/cloudness-io/cloudness/app/services/schema"
	specSvc "github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store/database"
	cliserver "github.com/cloudness-io/cloudness/cli/operations/server"
	"github.com/cloudness-io/cloudness/job"
	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/logstream"
	"github.com/cloudness-io/cloudness/pubsub"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

func initSystem(ctx context.Context, config *types.Config) (*cliserver.System, error) {
	wire.Build(
		cliserver.NewSystem,
		cliserver.ProvideDatabaseConfig,
		cliserver.ProvideLockConfig,
		cliserver.ProvidePubSubConfig,
		cliserver.ProvideRedis,
		bootstrap.WireSet,
		server.WireSet,
		router.WireSet,
		services.WireSet,
		job.WireSet,
		// openapi.WireSet,
		auth.WireSet,
		user.WireSet,
		dbtx.WireSet,
		authn.WireSet,
		instance.WireSet,
		serverCtrl.WireSet,
		tenant.WireSet,
		project.WireSet,
		favorite.WireSet,
		githubapp.WireSet,
		gitpublic.WireSet,
		database.WireSet,
		application.WireSet,
		variable.WireSet,
		volume.WireSet,
		environment.WireSet,
		deployment.WireSet,
		logs.WireSet,
		template.WireSet,
		sse.WireSet,
		logstream.WireSet,

		//services
		schema.WireSet,
		githubAppSvc.WireSet,
		gitpublicSvc.WireSet,
		managerSvc.WireSet,
		configSvc.WireSet,
		specSvc.WireSet,
		dnsSvc.WireSet,
		proxySvc.WireSet,
		cleanup.WireSet,

		//pipelinerm
		scheduler.WireSet,
		triggerer.WireSet,
		canceler.WireSet,
		manager.WireSet,
		client.WireSet,
		agent.WireSet,

		//commons
		lock.WireSet,
		pubsub.WireSet,
	)
	return &cliserver.System{}, nil
}
