package server

import (
	"context"
	"net/url"
	"strings"

	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type ServerGeneralUpdateModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServerNetworkUpdateModel struct {
	WildCardDomain  string           `json:"wildcard_domain"`
	IPV4            string           `json:"ipv4"`
	DNSProvider     enum.DNSProvider `json:"dns_provider"`
	DNSProviderAuth string           `json:"dns_provider_auth"`
}

type ServerBuilderUpdateModel struct {
	Enabled             bool    `json:"enabled,string"`
	IsBuildServer       bool    `json:"is_build_server,string"`
	MaxConcurrentBuilds int64   `json:"max_concurrent_builds,string"`
	MaxCPUPerBuild      float64 `json:"max_cpu_per_build,string"`
	MaxMemoryPerBuild   float64 `json:"max_memory_per_build,string"`
}

type ServerLimitsUpdateModel struct {
	SupportsOnlineExpansion bool  `json:"supports_online_expansion,string"`
	MinVolumeSize           int64 `json:"min_volume_size,string"`
}

func (c *Controller) UpdateGeneral(ctx context.Context, in *ServerGeneralUpdateModel) (*types.Server, error) {
	if err := c.sanitizeGeneralUpdateModel(in); err != nil {
		return nil, err
	}

	server, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	server.Name = in.Name
	server.Description = in.Description
	return c.serverStore.Update(ctx, server)
}

func (c *Controller) UpdateNetwork(ctx context.Context, in *ServerNetworkUpdateModel) (*types.Server, error) {
	if err := c.sanitizeNetworkUpdateModel(in); err != nil {
		return nil, err
	}

	//flags
	doProvisionSSL := false

	server, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	server.WildCardDomain = in.WildCardDomain
	server.IPV4 = in.IPV4
	server.DNSProvider = in.DNSProvider
	if server.DNSProvider != enum.DNSProviderNone {
		server.DNSProviderAuth = in.DNSProviderAuth
	} else {
		server.DNSProviderAuth = ""
	}

	if in.WildCardDomain != "" {
		wildcardDoamin, err := url.Parse(in.WildCardDomain)
		if err != nil {
			return nil, err
		}

		if instance.DNSValidationEnabled {
			wildcardHost := "*." + wildcardDoamin.Hostname()
			err = c.dnsSvc.ValidateHost(ctx, wildcardHost, server.IPV4, instance.DNSServers, in.DNSProvider)
			if err != nil {
				return nil, err
			}
		}

		if wildcardDoamin.Scheme == "https" {
			doProvisionSSL = true
			_, _, hostname := helpers.ParseFQDN(in.WildCardDomain)
			if err := c.proxySvc.ValidateToken(ctx, in.DNSProvider, in.DNSProviderAuth, hostname); err != nil {
				return nil, err
			}
		}
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		manager, err := c.factory.GetServerManager(server)
		if doProvisionSSL {
			if err != nil {
				return err
			}
			if err := manager.AddWildcardDomainWithSSL(ctx, server); err != nil {
				return err
			}
		} else {
			if err := manager.RemoveWildcardSSL(ctx, server); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error removing wildcard certificate")
			}
		}

		server, err = c.serverStore.Update(ctx, server)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (c *Controller) UpdateBuilder(ctx context.Context, in *ServerBuilderUpdateModel) (*types.Server, error) {
	server, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	if in.Enabled {
		server.BuildEnabled = in.Enabled
		server.IsBuildServer = in.IsBuildServer
		server.MaxConcurrentBuilds = in.MaxConcurrentBuilds
		server.MaxCPUPerBuild = in.MaxCPUPerBuild
		server.MaxMemoryPerBuild = in.MaxMemoryPerBuild
	} else {
		server.BuildEnabled = false
	}

	return c.serverStore.Update(ctx, server)
}

func (c *Controller) UpdateLimits(ctx context.Context, in *ServerLimitsUpdateModel) (*types.Server, error) {
	if in.MinVolumeSize <= 0 {
		errors := check.NewValidationErrors()
		errors.AddValidationError("min_volume_size", check.NewValidationError("minimum volume size must be greater than 0"))
		return nil, errors
	}

	server, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	server.VolumeMinSize = in.MinVolumeSize
	server.VolumeSupportsOnlineExpansion = in.SupportsOnlineExpansion

	return c.serverStore.Update(ctx, server)
}

func (c *Controller) sanitizeGeneralUpdateModel(in *ServerGeneralUpdateModel) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if err := check.Description(in.Description); err != nil {
		errors.AddValidationError("description", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) sanitizeNetworkUpdateModel(in *ServerNetworkUpdateModel) error {
	errors := check.NewValidationErrors()
	if err := check.FQDN(in.WildCardDomain); in.WildCardDomain != "" && err != nil {
		errors.AddValidationError("wildcard_domain", err)
	}
	if err := check.IPV4(in.IPV4); err != nil {
		errors.AddValidationError("ipv4", err)
	}

	if in.WildCardDomain != "" {
		if dnsProvider := enum.DNSProviderFromString(string(in.DNSProvider)); dnsProvider != "" {
			in.DNSProvider = dnsProvider
		} else {
			errors.AddValidationError("dns_provider", usererror.BadRequest("DNS Provider is not supported"))
		}

		if strings.HasPrefix(in.WildCardDomain, "https") {
			if in.DNSProvider == enum.DNSProviderNone {
				errors.AddValidationError("dns_provider", usererror.BadRequest("DNS Provider is required to provision SSL certificate for wildcard domains"))
			}
			if in.DNSProviderAuth == "" {
				errors.AddValidationError("dns_provider_auth", usererror.BadRequest("DNS Auth is required for https domain"))
			}
		}
	}

	if errors.HasError() {
		return errors
	}
	return nil
}
