package application

import (
	"context"
	"regexp"
	"strings"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

func (c *Controller) suggestPrivateDomain(ctx context.Context, app *types.Application) (string, error) {
	apps, err := c.List(ctx, app.TenantID, app.ProjectID, app.EnvironmentID)
	if err != nil {
		return "", err
	}

	privateSubDomain := helpers.Normalize(app.Name)

	for _, a := range apps {
		if a.PrivateDomain == privateSubDomain {
			return c.generateSubdomain(app), nil
		}
	}
	return privateSubDomain, nil
}

func (c *Controller) validateFQDN(ctx context.Context, fqdn string) error {
	errors := check.NewValidationErrors()
	if err := check.FQDN(fqdn); err != nil {
		errors.AddValidationError("fqdn", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) SuggestFQDN(ctx context.Context, app *types.Application) (string, error) {
	subDomain := c.generateSubdomain(app)

	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return "", err
	}

	domain, err := server.GetDomain()
	if err != nil {
		return "", err
	}

	return domain.Scheme + "://" + subDomain + "." + domain.Hostname, nil
}

func (c *Controller) generateSubdomain(app *types.Application) string {
	// Normalize name: lowercase, replace spaces with dashes
	normalized := helpers.Normalize(strings.ToLower(app.Name))
	normalized = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")

	sufix := helpers.GenerateSlug(8)

	return normalized + "-" + sufix

}
