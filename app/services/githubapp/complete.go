package githubapp

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/go-github/v69/github"
)

func (c *Service) CompleteManifest(ctx context.Context, ghApp *types.GithubApp, code string) error {
	ghClient, err := github.NewClient(nil).WithEnterpriseURLs(ghApp.ApiUrl, ghApp.ApiUrl)
	if err != nil {
		return err
	}

	app, _, err := ghClient.Apps.CompleteAppManifest(ctx, code)
	if err != nil {
		return err
	}

	now := time.Now().UTC().UnixMilli()
	ghApp.Name = app.GetName()
	ghApp.AppID = app.GetID()
	ghApp.ClientID = app.GetClientID()
	ghApp.ClientSecret = app.GetClientSecret()
	ghApp.WebhookSecret = app.GetWebhookSecret()
	ghApp.Updated = now

	privateKey := &types.PrivateKey{
		UID:      helpers.GenerateUID(),
		TenantID: ghApp.TenantID,
		Name:     ghApp.Name + "-private-key",
		Key:      app.GetPEM(),
		IsGit:    true,
		Created:  now,
		Updated:  now,
	}

	c.tx.WithTx(ctx, func(ctx context.Context) error {
		privateKey, err := c.privateKeyStore.Create(ctx, privateKey)
		if err != nil {
			return err
		}

		ghApp.PrivateKeyID = privateKey.ID
		_, err = c.githubAppStore.Update(ctx, ghApp)
		return err
	})

	return nil
}

func (c *Service) CompleteInstallation(ctx context.Context, ghApp *types.GithubApp, installationID int64) error {

	ghApp.InstallationID = installationID

	_, err := c.githubAppStore.Update(ctx, ghApp)
	if err != nil {
		return err
	}

	return nil
}
