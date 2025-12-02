package manager

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

type ServerManager interface {
	GetIP(ctx context.Context, server *types.Server) (string, error)
	//delets
	DeleteResources(ctx context.Context, server *types.Server, namespace string, identifier string) error
	DeleteNamespace(ctx context.Context, server *types.Server, namespace string) error
	DeleteVolume(ctx context.Context, server *types.Server, volume *types.Volume) error
	DeleteApplication(ctx context.Context, server *types.Server, app *types.Application) error

	//certificate related
	AddWildcardDomainWithSSL(ctx context.Context, server *types.Server) error
	RemoveWildcardSSL(ctx context.Context, server *types.Server) error
	AddSSLCertificate(ctx context.Context, server *types.Server, dns, certKey string, dnsProvider enum.DNSProvider, dnsAuthKey string) error
	RemoveSSLCertificate(ctx context.Context, server *types.Server, certKey string) error
	ListCertificates(ctx context.Context, server *types.Server) ([]*types.Certificate, error)

	AddHttpRoute(ctx context.Context, server *types.Server, namespace, name, service string, port int32, host string, httpScheme string) error
	RemoveHttpRoute(ctx context.Context, server *types.Server, namespace, name string) error

	//Registry endpoints
	AddOrUpdateRegistry(ctx context.Context, server *types.Server, name string, size int64, enableNodePort bool) error
	RemoveRegistry(ctx context.Context, server *types.Server, name string) error
	//Artifacts
	ListArtifacts(ctx context.Context, server *types.Server, app *types.Application) ([]*types.Artifact, error)
	TailLogs(ctx context.Context, server *types.Server, app *types.Application) (<-chan *types.ArtifactLogLine, <-chan error, error)
}
