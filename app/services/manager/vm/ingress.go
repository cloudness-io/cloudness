package vm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudness-io/cloudness/types"

	"gopkg.in/yaml.v3"
)

// traefikDynamicConfig represents a Traefik file provider dynamic configuration.
type traefikDynamicConfig struct {
	HTTP traefikHTTPConfig `yaml:"http"`
}

type traefikHTTPConfig struct {
	Routers  map[string]traefikRouter  `yaml:"routers"`
	Services map[string]traefikService `yaml:"services"`
}

type traefikRouter struct {
	Rule        string      `yaml:"rule"`
	EntryPoints []string    `yaml:"entryPoints"`
	Service     string      `yaml:"service"`
	TLS         *traefikTLS `yaml:"tls,omitempty"`
}

type traefikTLS struct {
	CertResolver string `yaml:"certResolver,omitempty"`
}

type traefikService struct {
	LoadBalancer traefikLoadBalancer `yaml:"loadBalancer"`
}

type traefikLoadBalancer struct {
	Servers []traefikLBServer `yaml:"servers"`
}

type traefikLBServer struct {
	URL string `yaml:"url"`
}

func httpRouteName(key string) string {
	return key + "-http-route"
}

func (m *VmManager) AddHttpRoute(ctx context.Context, server *types.Server, namespace, key, service string, port int32, host string, httpScheme string) error {
	routeName := httpRouteName(key)

	domain, err := server.GetDomain()
	if err != nil {
		return err
	}

	// If the host is under the wildcard domain and scheme is http,
	// use the plain "web" entrypoint. Otherwise use "websecure" with TLS.
	isWebListener := strings.HasSuffix(host, domain.Hostname) && httpScheme == "http"

	router := traefikRouter{
		Rule:    fmt.Sprintf("Host(`%s`)", host),
		Service: routeName,
	}

	if isWebListener {
		router.EntryPoints = []string{"web"}
	} else {
		router.EntryPoints = []string{"websecure"}
		router.TLS = &traefikTLS{
			CertResolver: "letsencrypt",
		}
	}

	// The service name is the container name on the Docker network (namespace).
	// Traefik, connected to the same network, resolves it by container name.
	cfg := traefikDynamicConfig{
		HTTP: traefikHTTPConfig{
			Routers: map[string]traefikRouter{
				routeName: router,
			},
			Services: map[string]traefikService{
				routeName: {
					LoadBalancer: traefikLoadBalancer{
						Servers: []traefikLBServer{
							{URL: fmt.Sprintf("http://%s:%d", service, port)},
						},
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal traefik config: %w", err)
	}

	configPath := filepath.Join(traefikDynamicConfigDir, routeName+".yml")

	if isLocal(server) {
		if err := os.MkdirAll(traefikDynamicConfigDir, 0o750); err != nil {
			return fmt.Errorf("failed to create traefik config directory: %w", err)
		}
		return os.WriteFile(configPath, data, 0o640)
	}

	// Remote: write via SSH
	cmd := fmt.Sprintf("mkdir -p %s && cat > %s << 'CLOUDNESS_EOF'\n%sCLOUDNESS_EOF",
		traefikDynamicConfigDir, configPath, string(data))
	_, err = m.runCmd(ctx, server, cmd)
	return err
}

func (m *VmManager) RemoveHttpRoute(ctx context.Context, server *types.Server, _, key string) error {
	configPath := filepath.Join(traefikDynamicConfigDir, httpRouteName(key)+".yml")

	if isLocal(server) {
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove traefik config: %w", err)
		}
		return nil
	}

	cmd := fmt.Sprintf("rm -f %s", configPath)
	_, err := m.runCmd(ctx, server, cmd)
	return err
}
