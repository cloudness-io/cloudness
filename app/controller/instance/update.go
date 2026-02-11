package instance

import (
	"context"
	"net/url"
	"strings"

	"github.com/cloudness-io/cloudness/app/services/manager/kube"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type InstanceUpdateFQDNInput struct {
	FQDN            string           `json:"fqdn"`
	DNSProvider     enum.DNSProvider `json:"dns_provider"`
	DNSProviderAuth string           `json:"dns_provider_auth"`
}

type InstanceUpdateDNSConfigInput struct {
	DNSValidationEnabled bool   `json:"dns_validation_enabled,string"`
	DNSServers           string `json:"dns_servers"`
}

type InstanceUpdateScriptsInput struct {
	ExternalScripts   string `json:"external_scripts"`
	AdditionalScripts string `json:"additional_scripts"`
}

type InstanceUpdateInput struct {
	FQDN                 string           `json:"fqdn"`
	DNSProvider          enum.DNSProvider `json:"dns_provider"`
	DNSProviderAuth      string           `json:"dns_provider_auth"`
	DNSValidationEnabled bool             `json:"dns_validation_enabled,string"`
	DNSServers           string           `json:"dns_servers"`
	ExternalScripts      string           `json:"external_scripts"`
	AdditionalScripts    string           `json:"additional_scripts"`
}

type InstanceRegistryUpdateInput struct {
	Enabled       bool  `json:"registry_enabled,string"`
	Size          int64 `json:"registry_size,string"`
	MirrorEnabled bool  `json:"registry_mirror_enabled,string"`
	MirrorSize    int64 `json:"registry_mirror_size,string"`
}

func (c *Controller) UpdateFQDN(ctx context.Context, server *types.Server, in *InstanceUpdateFQDNInput) (*types.Instance, error) {
	if err := c.sanitizeFQDNUpdateModel(in); err != nil {
		return nil, err
	}

	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	//flags
	hadFQDN := false
	doAddRoute := false
	doProvisionSSL := false

	fqdnURL, err := url.Parse(in.FQDN)
	if err != nil {
		return nil, err
	}

	if in.FQDN == "" && instance.FQDN != "" {
		hadFQDN = true
	}

	if in.FQDN != "" && instance.FQDN != in.FQDN {
		doAddRoute = true
		isServerSubdomain := false
		instanceSchema, _, instanceDomain := helpers.ParseFQDN(in.FQDN)
		if server.WildCardDomain != "" {
			//validate if instance is subdomain of server
			_, _, serverDomain := helpers.ParseFQDN(server.WildCardDomain)
			if strings.HasSuffix(fqdnURL.Hostname(), serverDomain) {
				isServerSubdomain = true
			}
		}

		log.Ctx(ctx).Debug().
			Str("instanceSchema", instanceSchema).
			Str("instanceDomain", instanceDomain).
			Bool("isServerSuddomain", isServerSubdomain).
			Msg("Control flags for server subdomain")

		if !isServerSubdomain {
			// validate dns record
			if instance.DNSValidationEnabled {
				err = c.dnsSvc.ValidateHost(ctx, fqdnURL.Hostname(), server.IPV4, instance.DNSServers, in.DNSProvider)
				if err != nil {
					return nil, err
				}
			}

			if instanceSchema == "https" {
				doProvisionSSL = true
				if err := c.proxySvc.ValidateToken(ctx, in.DNSProvider, in.DNSProviderAuth, instanceDomain); err != nil {
					return nil, err
				}
			}
		}
	}

	log.Ctx(ctx).Debug().Str("Hostname", fqdnURL.Hostname()).Bool("doProvisionSSL", doProvisionSSL).Bool("doAddRoute", doAddRoute).Bool("hadFQDN", hadFQDN).Msg("Control flags for domain")

	instance.FQDN = in.FQDN
	if instance.FQDN == "" {
		instance.DNSProvider = enum.DNSProviderNone
		instance.DNSProviderAuth = ""
	} else {
		instance.DNSProvider = in.DNSProvider
		instance.DNSProviderAuth = in.DNSProviderAuth
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		manager, err := c.factory.GetServerManager(server)
		if err != nil {
			return err
		}
		certKey := "cloudness-app-route"
		if doProvisionSSL {
			if err := manager.AddSSLCertificate(ctx, server, kube.DefaultK8sCloudnessName, fqdnURL.Hostname(), certKey, instance.DNSProvider, instance.DNSProviderAuth); err != nil {
				return err
			}
		}

		if doAddRoute {
			err := manager.AddHttpRoute(ctx, server, kube.DefaultK8sCloudnessNamespace, certKey, kube.DefaultK8sCloudnessService, kube.DefaultK8sCloudnessPort, fqdnURL.Hostname(), fqdnURL.Scheme)
			if err != nil {
				return err
			}
		} else if hadFQDN {
			log.Ctx(ctx).Debug().Msg("removing http route")
			if err := manager.RemoveHttpRoute(ctx, server, kube.DefaultK8sCloudnessNamespace, certKey); err != nil {
				return err
			}
		}

		instance, err = c.instanceStore.Update(ctx, instance)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Controller) UpdateDNSConfig(ctx context.Context, in *InstanceUpdateDNSConfigInput) (*types.Instance, error) {
	if err := c.sanitizeDNSConfigUpdateModel(in); err != nil {
		return nil, err
	}
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	instance.DNSValidationEnabled = in.DNSValidationEnabled
	instance.DNSServers = in.DNSServers

	instance, err = c.instanceStore.Update(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Controller) UpdateScripts(ctx context.Context, in *InstanceUpdateScriptsInput) (*types.Instance, error) {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	instance.ExternalScripts = in.ExternalScripts
	instance.AdditionalScripts = in.AdditionalScripts

	instance, err = c.instanceStore.Update(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Controller) UpdateEnableSignup(ctx context.Context, enableSignup bool) error {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return err
	}

	instance.UserSignupEnabled = enableSignup
	_, err = c.instanceStore.Update(ctx, instance)
	return err
}

func (c *Controller) UpdateDemoUser(ctx context.Context, demoUserEnabled bool) error {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return err
	}

	instance.DemoUserEnabled = demoUserEnabled
	_, err = c.instanceStore.Update(ctx, instance)
	return err
}

func (c *Controller) UpdateSuperAdmin(ctx context.Context, superAdmin *int64) error {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return err
	}

	instance.SuperAdmin = superAdmin
	_, err = c.instanceStore.Update(ctx, instance)
	return err
}

func (c *Controller) UpdateRegistry(ctx context.Context, in *InstanceRegistryUpdateInput) (*types.Instance, error) {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}
	// validations
	errors := check.NewValidationErrors()
	if in.Enabled {
		if in.Size < instance.RegistrySize {
			errors.AddValidationError("registry_size", check.NewValidationError("Volume size cannot be less than current size"))
		}
		instance.RegistrySize = in.Size
	} else {
		errors.AddValidationError("registry_enabled", check.NewValidationError("Registry cannot be disabled"))
	}
	if in.MirrorEnabled {
		if instance.RegistryMirrorEnabled == true && in.MirrorSize < instance.RegistryMirrorSize {
			errors.AddValidationError("registry_proxy_size", check.NewValidationError("Volume size cannot be less than current size"))
		}
		instance.RegistryMirrorSize = in.MirrorSize
	}
	if errors.HasError() {
		return nil, errors
	}

	instance.RegistryEnabled = in.Enabled
	instance.RegistryMirrorEnabled = in.MirrorEnabled

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		if _, err := c.instanceStore.Update(ctx, instance); err != nil {
			return err
		}
		server, err := c.serverStore.Find(ctx, 1)
		if err != nil {
			return err
		}

		factory, err := c.factory.GetServerManager(server)
		if err != nil {
			return err
		}

		if instance.RegistryEnabled {
			if err = factory.AddOrUpdateRegistry(ctx, server, registryKey, instance.RegistrySize, true); err != nil {
				return err
			}
		} else {
			if err = factory.RemoveRegistry(ctx, server, registryKey); err != nil {
				return err
			}
		}

		if instance.RegistryMirrorEnabled {
			if err = factory.AddOrUpdateRegistry(ctx, server, registryMirrorKey, instance.RegistryMirrorSize, false); err != nil {
				return err
			}
		} else {
			if err = factory.RemoveRegistry(ctx, server, registryMirrorKey); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
	// return c.instanceStore.Update(ctx, instance)
}

func (c *Controller) sanitizeFQDNUpdateModel(in *InstanceUpdateFQDNInput) error {
	errors := check.NewValidationErrors()
	if err := check.FQDN(in.FQDN); in.FQDN != "" && err != nil {
		errors.AddValidationError("fqdn", err)
	}

	if in.FQDN != "" {
		if dnsProvider := enum.DNSProviderFromString(string(in.DNSProvider)); dnsProvider != "" {
			in.DNSProvider = dnsProvider
		} else {
			errors.AddValidationError("dns_provider", usererror.BadRequest("DNS Provider is not supported"))
		}

		if in.DNSProvider != enum.DNSProviderNone && strings.HasPrefix(in.FQDN, "https") {
			if in.DNSProviderAuth == "" {
				errors.AddValidationError("dns_provider_auth", usererror.BadRequest("DNS Auth is required for https domain behind a DNS proxy provider"))
			}
		}
	}

	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) sanitizeDNSConfigUpdateModel(in *InstanceUpdateDNSConfigInput) error {
	errors := check.NewValidationErrors()

	if in.DNSValidationEnabled {
		if in.DNSServers == "" {
			errors.AddValidationError("dns_servers", usererror.BadRequest("DNS servers are required"))
		}
	}

	if errors.HasError() {
		return errors
	}
	return nil
}
