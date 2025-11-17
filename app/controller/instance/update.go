package instance

import (
	"context"
	"net/url"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type InstanceUpdateInput struct {
	FQDN                 string `json:"fqdn"`
	DNSValidationEnabled bool   `json:"dns_validation_enabled,string"`
	DNSServers           string `json:"dns_servers"`
	ExternalScripts      string `json:"external_scripts"`
}

type InstanceRegistryUpdateInput struct {
	Enabled       bool  `json:"registry_enabled,string"`
	Size          int64 `json:"registry_size,string"`
	MirrorEnabled bool  `json:"registry_mirror_enabled,string"`
	MirrorSize    int64 `json:"registry_mirror_size,string"`
}

func (c *Controller) Update(ctx context.Context, server *types.Server, in *InstanceUpdateInput) (*types.Instance, error) {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	if in.DNSValidationEnabled && instance.FQDN != in.FQDN {
		fqdnURL, err := url.Parse(in.FQDN)
		if err != nil {
			return nil, err
		}

		err = c.dnsSvc.ValidateHost(ctx, fqdnURL.Hostname(), server.IPV4, in.DNSServers, server.DNSProxy)
		if err != nil {
			return nil, err
		}
	}

	instance.FQDN = in.FQDN
	instance.DNSValidationEnabled = in.DNSValidationEnabled
	instance.DNSServers = in.DNSServers
	instance.ExternalScripts = in.ExternalScripts

	return c.instanceStore.Update(ctx, instance)
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
