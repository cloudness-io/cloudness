package manager

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

type ServerManager interface {
	GetIP(ctx context.Context, server *types.Server) (string, error)
	DeleteResources(ctx context.Context, server *types.Server, namespace string, identifier string) error
	DeleteNamespace(ctx context.Context, server *types.Server, namespace string) error
	DeleteVolume(ctx context.Context, server *types.Server, volume *types.Volume) error
	DeleteApplication(ctx context.Context, server *types.Server, app *types.Application) error
	AddWildcardDomainWithSSL(ctx context.Context, server *types.Server) error
	ListCertificates(ctx context.Context, server *types.Server) ([]*types.Certificate, error)
	//Registry endpoints
	AddOrUpdateRegistry(ctx context.Context, server *types.Server, name string, size int64, enableNodePort bool) error
	RemoveRegistry(ctx context.Context, server *types.Server, name string) error
	//Artifacts
	ListArtifacts(ctx context.Context, server *types.Server, app *types.Application) ([]*types.Artifact, error)
	TailLogs(ctx context.Context, server *types.Server, app *types.Application) (<-chan *types.ArtifactLogLine, <-chan error, error)
}
