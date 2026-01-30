package main

import (
	"github.com/cloudness-io/cloudness/cli"
	"github.com/cloudness-io/cloudness/cli/operations/server"

	"github.com/alecthomas/kingpin/v2"
)

const (
	applicationName = "Cloudness Platform"
	description     = "Awesome App Platform"
)

func main() {
	args := cli.GetArguments()

	app := kingpin.New(applicationName, description)

	server.Register(app, initSystem)

	// swagger.Register(app, openapi.NewOpenAPIService())

	kingpin.MustParse(app.Parse(args))

}
