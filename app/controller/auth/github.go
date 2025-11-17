package auth

import (
	"context"
	"encoding/json"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func (c *Controller) githubOAuth2Config(authSetting *types.AuthSetting) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     authSetting.ClientID,
		ClientSecret: authSetting.ClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
}

func (c *Controller) githubCallback(ctx context.Context, authSetting *types.AuthSetting, code string) (string, error) {
	config := c.githubOAuth2Config(authSetting)

	token, err := config.Exchange(ctx, code, oauth2.AccessTypeOnline)
	if err != nil {
		return "", err
	}

	client := config.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Primary  bool   `json:"primary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	var userEmail string
	for _, e := range emails {
		if e.Primary {
			userEmail = e.Email
			break
		}
	}

	if userEmail == "" {
		return "", errors.BadRequest("No primary email id found")
	}

	return userEmail, nil
}
