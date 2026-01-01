package main

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Deployer handles Kubernetes deployments using hybrid approach:
// - kubectl for applying manifests (battle-tested)
// - client-go for watching/monitoring (efficient)
type Deployer struct {
	config   *Config
	Operator Operator
	log      *Logger
}

// NewDeployer creates a new deployer instance
func NewDeployer(cfg *Config) (*Deployer, error) {
	// Get in-cluster config
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	// Create clientset for watching
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	logger := NewLogger(cfg.Verbose)
	var operator Operator
	switch cfg.AppType {
	case AppTypeStateless:
		operator = NewStatelessOperator(cfg, clientset, logger)
	case AppTypeStateful:
		operator = NewStatefulOperator(cfg, clientset, logger)
	}

	return &Deployer{
		config:   cfg,
		Operator: operator,
		log:      logger,
	}, nil
}

// Deploy runs the full deployment workflow
func (k *Deployer) Deploy(ctx context.Context) error {
	k.log.Section("Deploying application")

	// Step 1: Apply common artifacts (namespace, service account, etc.)
	if err := k.Operator.ApplyCommon(ctx); err != nil {
		return err
	}
	k.log.Step("Access control and storage artifacts configured")

	// Step 2: Handle volumes (if any)
	if err := k.Operator.Volumes(ctx); err != nil {
		return err
	}

	// Step 3: Deploy application
	if err := k.Operator.Deploy(ctx); err != nil {
		return err
	}
	k.log.Step("Application deployed")

	// Step 4: Configure routes (if any)
	if err := k.Operator.Ingress(ctx); err != nil {
		return err
	}

	// Cleanup: Remove opposite resource type
	k.Operator.Cleanup(ctx)

	k.log.Success("Deployment completed successfully!")
	return nil
}
