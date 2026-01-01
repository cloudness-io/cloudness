package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all deployment configuration from environment variables
type Config struct {
	// Application identifiers
	AppIdentifier string
	AppNamespace  string
	AppType       AppType // "Stateless" or "Stateful"

	// Feature flags
	HasVolume   bool
	HasRoute    bool
	NeedRemount bool

	// YAML file paths
	DeployPath     string
	CommonYAMLPath string
	VolumeYAMLPath string
	AppYAMLPath    string
	RouteYAMLPath  string

	// Timeouts
	RolloutTimeoutStateless int // seconds
	RolloutTimeoutStateful  int // seconds
	PVCResizeTimeout        int // seconds
	PVCResizePollInterval   int // seconds

	// Options
	Verbose bool
}

type AppType string

const (
	AppTypeStateless AppType = "Stateless"
	AppTypeStateful  AppType = "Stateful"
)

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		AppIdentifier: os.Getenv("CLOUDNESS_DEPLOY_APP_IDENTIFIER"),
		AppNamespace:  os.Getenv("CLOUDNESS_DEPLOY_APP_NAMESPACE"),
		AppType:       AppType(os.Getenv("CLOUDNESS_DEPLOY_FLAG_APP_TYPE")),
		DeployPath:    os.Getenv("CLOUDNESS_DEPLOY_PATH"),
		HasVolume:     os.Getenv("CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME") == "1",
		HasRoute:      os.Getenv("CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE") == "1",
		NeedRemount:   os.Getenv("CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT") == "1",
		Verbose:       os.Getenv("VERBOSE") == "true",

		// Defaults
		RolloutTimeoutStateless: getEnvInt("ROLLOUT_TIMEOUT_STATELESS", 60),
		RolloutTimeoutStateful:  getEnvInt("ROLLOUT_TIMEOUT_STATEFUL", 120),
		PVCResizeTimeout:        getEnvInt("PVC_RESIZE_TIMEOUT", 300),
		PVCResizePollInterval:   getEnvInt("PVC_RESIZE_POLL_INTERVAL", 5),
	}

	// Validate required fields
	if cfg.AppIdentifier == "" {
		return nil, fmt.Errorf("CLOUDNESS_DEPLOY_APP_IDENTIFIER is required")
	}
	if cfg.AppNamespace == "" {
		return nil, fmt.Errorf("CLOUDNESS_DEPLOY_APP_NAMESPACE is required")
	}
	if cfg.AppType != AppTypeStateless && cfg.AppType != AppTypeStateful {
		return nil, fmt.Errorf("CLOUDNESS_DEPLOY_FLAG_APP_TYPE must be 'Stateless' or 'Stateful', got '%s'", cfg.AppType)
	}
	if cfg.DeployPath == "" {
		return nil, fmt.Errorf("CLOUDNESS_DEPLOY_PATH is required")
	}

	// Set YAML file paths
	cfg.CommonYAMLPath = cfg.DeployPath + "/common.yaml"
	cfg.VolumeYAMLPath = cfg.DeployPath + "/volume.yaml"
	cfg.AppYAMLPath = cfg.DeployPath + "/app.yaml"
	cfg.RouteYAMLPath = cfg.DeployPath + "/route.yaml"

	return cfg, nil
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

// ResourceType returns the Kubernetes resource type for this app
func (c *Config) ResourceType() string {
	if c.AppType == AppTypeStateless {
		return "Deployment"
	}
	return "StatefulSet"
}

// OppositeResourceType returns the opposite resource type (for cleanup)
func (c *Config) OppositeResourceType() string {
	if c.AppType == AppTypeStateless {
		return "StatefulSet"
	}
	return "Deployment"
}

// RolloutTimeout returns the appropriate timeout based on app type
func (c *Config) RolloutTimeout() int {
	if c.AppType == AppTypeStateless {
		return c.RolloutTimeoutStateless
	}
	return c.RolloutTimeoutStateful
}
