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

type InstanceUpdateInput struct {
	FQDN                 string           `json:"fqdn"`
	DNSProvider          enum.DNSProvider `json:"dns_provider"`
	DNSProviderAuth      string           `json:"dns_provider_auth"`
	DNSValidationEnabled bool             `json:"dns_validation_enabled,string"`
	DNSServers           string           `json:"dns_servers"`
	ExternalScripts      string           `json:"external_scripts"`
}

type InstanceRegistryUpdateInput struct {
	Enabled       bool  `json:"registry_enabled,string"`
	Size          int64 `json:"registry_size,string"`
	MirrorEnabled bool  `json:"registry_mirror_enabled,string"`
	MirrorSize    int64 `json:"registry_mirror_size,string"`
}

func (c *Controller) Update(ctx context.Context, server *types.Server, in *InstanceUpdateInput) (*types.Instance, error) {
	if err := c.sanitizeGeneralUpdateModel(in); err != nil {
		return nil, err
	}
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	//flags
	doAddRoute := false
	doProvisionSSL := false

	fqdnURL, err := url.Parse(in.FQDN)
	// if instance.FQDN != in.FQDN && in.FQDN != "" {
	if in.FQDN != "" {
		doAddRoute = true
		if err != nil {
			return nil, err
		}

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
			if in.DNSValidationEnabled {
				err = c.dnsSvc.ValidateHost(ctx, fqdnURL.Hostname(), server.IPV4, in.DNSServers, in.DNSProvider)
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

	log.Ctx(ctx).Debug().Str("Hostname", fqdnURL.Hostname()).Bool("doProvisionSSL", doProvisionSSL).Bool("doAddRoute", doAddRoute).Msg("Control flags for domain")

	instance.FQDN = in.FQDN
	if instance.FQDN == "" {
		instance.DNSProvider = enum.DNSProviderNone
		instance.DNSProviderAuth = ""
	} else {
		instance.DNSProvider = in.DNSProvider
		instance.DNSProviderAuth = in.DNSProviderAuth
	}
	instance.DNSValidationEnabled = in.DNSValidationEnabled
	instance.DNSServers = in.DNSServers
	instance.ExternalScripts = in.ExternalScripts
	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		manager, err := c.factory.GetServerManager(server)
		if err != nil {
			return err
		}
		if doProvisionSSL {
			if err := manager.AddSSLCertificate(ctx, server, kube.DefaultK8sCloudnessNamespace, fqdnURL.Hostname(), "cloudness-app-certificate", instance.DNSProvider, instance.DNSProviderAuth); err != nil {
				return err
			}
		} else {
			if err := manager.RemoveSSLCertificate(ctx, server, kube.DefaultK8sCloudnessNamespace, "cloudness-app-certificate"); err != nil {
				return err
			}
		}

		if doAddRoute {
			if err := manager.AddHttpRoute(ctx, server, kube.DefaultK8sCloudnessNamespace, "cloudness-custom-http", kube.DefaultK8sCloudnessService, kube.DefaultK8sCloudnessPort, fqdnURL.Hostname()); err != nil {
				return err
			}
		} else {
			if err := manager.RemoveHttpRoute(ctx, server, kube.DefaultK8sCloudnessNamespace, "cloudness-custom-http"); err != nil {
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

func (c *Controller) UpdateWithServer(ctx context.Context, server *types.Server) (*types.Instance, error) {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	instance.PublicIPV4 = server.IPV4
	instance.PublicIPV6 = server.IPV6

	return c.instanceStore.Update(ctx, instance)
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

func (c *Controller) sanitizeGeneralUpdateModel(in *InstanceUpdateInput) error {
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
