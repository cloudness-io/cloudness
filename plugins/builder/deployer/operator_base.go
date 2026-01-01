package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
)

type BaseOpeator struct {
	config  *Config
	kubectl *Kubectl
	log     *Logger
}

func NewBaseOperator(config *Config, kubectl *Kubectl, log *Logger) *BaseOpeator {
	return &BaseOpeator{
		config:  config,
		kubectl: kubectl,
		log:     log,
	}
}

func (b *BaseOpeator) ApplyCommon(ctx context.Context) error {
	if err := b.kubectl.ApplyYAMLFile(ctx, b.config.CommonYAMLPath); err != nil {
		return fmt.Errorf("failed to apply common artifacts: %w", err)
	}
	return nil
}

func (b *BaseOpeator) ApplyIngress(ctx context.Context) error {
	if b.config.HasRoute {
		if err := b.kubectl.ApplyYAMLFile(ctx, b.config.RouteYAMLPath); err != nil {
			return fmt.Errorf("failed to deploy routes: %w", err)
		}
		b.log.Step("HTTP routes configured")
	}
	return nil
}

// utils
func isDeploymentReady(deploy *appsv1.Deployment) bool {
	if deploy.Generation != deploy.Status.ObservedGeneration {
		return false
	}

	if deploy.Spec.Replicas == nil {
		return false
	}

	replicas := *deploy.Spec.Replicas
	return deploy.Status.UpdatedReplicas == replicas &&
		deploy.Status.ReadyReplicas == replicas &&
		deploy.Status.AvailableReplicas == replicas
}

func isStatefulSetReady(sts *appsv1.StatefulSet) bool {
	if sts.Generation != sts.Status.ObservedGeneration {
		return false
	}

	if sts.Spec.Replicas == nil {
		return false
	}

	replicas := *sts.Spec.Replicas
	return sts.Status.UpdatedReplicas == replicas &&
		sts.Status.ReadyReplicas == replicas
}
