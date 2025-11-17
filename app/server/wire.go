package server

import (
	"github.com/cloudness-io/cloudness/app/router"
	"github.com/cloudness-io/cloudness/http"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(ProvideServer)

// ProvideServer provides a server instance.
func ProvideServer(config *types.Config, router *router.Router) *Server {
	return &Server{
		http.NewServer(
			http.Config{
				Port:     config.Server.HTTP.Port,
				Acme:     config.Server.Acme.Enabled,
				AcmeHost: config.Server.Acme.Host,
			},
			router,
		),
	}
}
