package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"gopkg.in/yaml.v3"
)

// traefikStaticConfig fragments used when configuring ACME on the Traefik instance.
type traefikCertConfig struct {
	HTTP traefikCertHTTPConfig `yaml:"http,omitempty"`
	TLS  *traefikTLSStore      `yaml:"tls,omitempty"`
}

type traefikCertHTTPConfig struct {
	Routers map[string]traefikRouter `yaml:"routers,omitempty"`
}

type traefikTLSStore struct {
	Stores map[string]traefikTLSStoreEntry `yaml:"stores,omitempty"`
}

type traefikTLSStoreEntry struct {
	DefaultCertificate *traefikDefaultCert `yaml:"defaultCertificate,omitempty"`
}

type traefikDefaultCert struct {
	CertFile string `yaml:"certFile,omitempty"`
	KeyFile  string `yaml:"keyFile,omitempty"`
}

// certFileName returns the config file name for a given certificate key.
func certFileName(certKey string) string {
	return certKey + "-cert.yml"
}

func (m *VmManager) AddWildcardDomainWithSSL(ctx context.Context, server *types.Server) error {
	wURL, err := url.Parse(server.WildCardDomain)
	if err != nil {
		return nil
	}

	if wURL.Scheme == "https" {
		wildcardDNS := fmt.Sprintf("*.%s", wURL.Hostname())
		return m.AddSSLCertificate(ctx, server, "", wildcardDNS, "cloudness-wildcard", server.DNSProvider, server.DNSProviderAuth)
	}

	return nil
}

func (m *VmManager) RemoveWildcardSSL(ctx context.Context, server *types.Server) error {
	return m.RemoveSSLCertificate(ctx, server, "", "cloudness-wildcard")
}

func (m *VmManager) AddSSLCertificate(ctx context.Context, server *types.Server, namespace, dns, certKey string, dnsProvider enum.DNSProvider, dnsAuthKey string) error {
	// Traefik handles ACME certificate provisioning natively.
	// We create a dynamic config that:
	// 1. Adds a router for the ACME HTTP-01 challenge (if not using DNS challenge)
	// 2. Sets up TLS via the Traefik entrypoint (letsencrypt certResolver)
	//
	// For DNS challenges (e.g. Cloudflare), the Traefik static config must
	// already include the ACME DNS provider settings. The dynamic config just
	// ensures routing exists for the domain.

	isWildcard := strings.HasPrefix(dns, "*.")
	routerName := certKey + "-acme"

	router := traefikRouter{
		Rule:        fmt.Sprintf("Host(`%s`)", dns),
		EntryPoints: []string{"websecure"},
		Service:     "noop@internal",
		TLS: &traefikTLS{
			CertResolver: "letsencrypt",
		},
	}

	// For wildcard domains with DNS challenge, we rely on the Traefik static config
	// having the ACME dnsChallenge section configured (e.g. for Cloudflare).
	// The dynamic config router triggers certificate issuance.
	if isWildcard && dnsProvider == enum.DNSProviderCloudflare && dnsAuthKey != "" {
		// Write the Cloudflare env vars so Traefik can use them for DNS challenge.
		envCmd := fmt.Sprintf(
			"docker exec cloudness-proxy sh -c 'echo CF_DNS_API_TOKEN=%s > /tmp/.cf-env' 2>/dev/null || true",
			dnsAuthKey,
		)
		if _, err := m.runCmd(ctx, server, envCmd); err != nil {
			// Non-fatal: Traefik may already have static env config
			_ = err
		}
	}

	cfg := traefikCertConfig{
		HTTP: traefikCertHTTPConfig{
			Routers: map[string]traefikRouter{
				routerName: router,
			},
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal cert config: %w", err)
	}

	configPath := filepath.Join(traefikDynamicConfigDir, certFileName(certKey))

	if isLocal(server) {
		if err := os.MkdirAll(traefikDynamicConfigDir, 0o750); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		return os.WriteFile(configPath, data, 0o640)
	}

	// Remote: write via SSH
	cmd := fmt.Sprintf("mkdir -p %s && cat > %s << 'CLOUDNESS_EOF'\n%sCLOUDNESS_EOF",
		traefikDynamicConfigDir, configPath, string(data))
	_, err = m.runCmd(ctx, server, cmd)
	return err
}

func (m *VmManager) RemoveSSLCertificate(ctx context.Context, server *types.Server, namespace, certKey string) error {
	configPath := filepath.Join(traefikDynamicConfigDir, certFileName(certKey))

	if isLocal(server) {
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove cert config: %w", err)
		}
		return nil
	}

	cmd := fmt.Sprintf("rm -f %s", configPath)
	_, err := m.runCmd(ctx, server, cmd)
	return err
}

// acmeCertInfo represents a certificate entry in Traefik's acme.json.
type acmeCertInfo struct {
	Domain      acmeDomain `json:"domain"`
	Certificate string     `json:"certificate"`
	Key         string     `json:"key"`
}

type acmeDomain struct {
	Main string   `json:"main"`
	SANs []string `json:"sans"`
}

type acmeAccount struct {
	Certificates []acmeCertInfo `json:"Certificates"`
}

func (m *VmManager) ListCertificates(ctx context.Context, server *types.Server) ([]*types.Certificate, error) {
	// Traefik stores ACME certificates in acme.json.
	cmd := fmt.Sprintf("cat %s 2>/dev/null || echo '{}'", traefikACMEStoragePath)
	out, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to read acme.json: %w", err)
	}

	// acme.json structure: { "letsencrypt": { "Certificates": [...] } }
	var acmeData map[string]*acmeAccount
	if err := json.Unmarshal([]byte(strings.TrimSpace(out.Stdout)), &acmeData); err != nil {
		return nil, fmt.Errorf("failed to parse acme.json: %w", err)
	}

	now := time.Now().UTC()
	var certs []*types.Certificate

	for resolverName, account := range acmeData {
		if account == nil {
			continue
		}
		for _, cert := range account.Certificates {
			dnsNames := []string{cert.Domain.Main}
			dnsNames = append(dnsNames, cert.Domain.SANs...)

			certs = append(certs, &types.Certificate{
				Name:      cert.Domain.Main,
				Namespace: resolverName,
				DNSNames:  dnsNames,
				IssuerRef: "ACME/" + resolverName,
				Ready:     "True",
				NotBefore: now,
			})
		}
	}

	return certs, nil
}
