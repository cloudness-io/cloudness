package schema

import "github.com/google/wire"

var WireSet = wire.NewSet(
	ProviderSchemaService,
)

func ProviderSchemaService() *Service {
	return NewService()
}
