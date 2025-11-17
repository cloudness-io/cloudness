package server

import (
	"github.com/cloudness-io/cloudness/app/bootstrap"
	"github.com/cloudness-io/cloudness/app/pipeline/agent"
	"github.com/cloudness-io/cloudness/app/server"
	"github.com/cloudness-io/cloudness/app/services"
)

// System stores high level System sub-routines.
type System struct {
	bootstrap bootstrap.Bootstrap
	server    *server.Server
	agent     *agent.Agent
	services  services.Services
}

// NewSystem returns a new system structure.
func NewSystem(
	bootstrap bootstrap.Bootstrap,
	server *server.Server,
	agent *agent.Agent,
	services services.Services,
) *System {
	return &System{
		bootstrap: bootstrap,
		server:    server,
		agent:     agent,
		services:  services,
	}
}
