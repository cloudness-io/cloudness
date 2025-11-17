package githubapp

import (
	"context"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

func (s *Service) getGithubClient(ctx context.Context, ghApp *types.GithubApp) (*github.Client, error) {

	accessToken, err := s.generateAccessToken(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	client := oauth2.NewClient(ctx, tokenSource)
	ghClient := github.NewClient(client)
	return ghClient, nil
}

func (s *Service) generateAccessToken(ctx context.Context, ghApp *types.GithubApp) (string, error) {
	jwtToken, err := s.generateJWT(ctx, ghApp)
	if err != nil {
		return "", err
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: jwtToken})
	client := oauth2.NewClient(ctx, tokenSource)

	ghClient, err := github.NewClient(client).WithEnterpriseURLs(ghApp.ApiUrl, ghApp.ApiUrl)
	if err != nil {
		return "", err
	}

	installation, _, err := ghClient.Apps.CreateInstallationToken(ctx, ghApp.InstallationID, &github.InstallationTokenOptions{})
	if err != nil {
		return "", err
	}

	return installation.GetToken(), nil
}

func (s *Service) generateJWT(ctx context.Context, ghApp *types.GithubApp) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("%d", ghApp.AppID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 10)),
	})

	key, err := s.loadPrivateKey(ctx, ghApp)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

func (s *Service) loadPrivateKey(ctx context.Context, ghApp *types.GithubApp) (*rsa.PrivateKey, error) {
	privateKey, err := s.privateKeyStore.Find(ctx, ghApp.TenantID, ghApp.PrivateKeyID)
	if err != nil {
		return nil, err
	}
	keyBytes := []byte(privateKey.Key)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
}
