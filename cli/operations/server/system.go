package server

import (
	"github.com/cloudness-io/cloudness/app/bootstrap"
	"github.com/cloudness-io/cloudness/app/pipeline/agent"
	"github.com/cloudness-io/cloudness/app/server"
	"github.com/cloudness-io/cloudness/app/services"
	"github.com/cloudness-io/cloudness/app/services/background"
)

// System stores high level System sub-routines.
type System struct {
	bootstrap  bootstrap.Bootstrap
	server     *server.Server
	agent      *agent.Agent
	services   services.Services
	background *background.Service
}

// NewSystem returns a new system structure.
func NewSystem(
	bootstrap bootstrap.Bootstrap,
	server *server.Server,
	agent *agent.Agent,
	services services.Services,
	background *background.Service,
) *System {
	return &System{
		bootstrap:  bootstrap,
		server:     server,
		agent:      agent,
		services:   services,
		background: background,
	}
}
