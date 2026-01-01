package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type StatelessOperator struct {
	config    *Config
	base      *BaseOpeator
	clientset kubernetes.Interface
	kubectl   *Kubectl
	log       *Logger
}

func NewStatelessOperator(cfg *Config, clientset kubernetes.Interface, log *Logger) *StatelessOperator {
	kubectl := NewKubectl(cfg, log)
	baseOperator := NewBaseOperator(cfg, kubectl, log)
	return &StatelessOperator{
		config:    cfg,
		base:      baseOperator,
		clientset: clientset,
		kubectl:   kubectl,
		log:       log,
	}
}

func (s *StatelessOperator) ApplyCommon(ctx context.Context) error {
	return s.base.ApplyCommon(ctx)
}

func (s *StatelessOperator) Volumes(ctx context.Context) error {
	return nil
}

func (s *StatelessOperator) Deploy(ctx context.Context) error {
	if err := s.kubectl.ApplyYAMLFile(ctx, s.config.AppYAMLPath); err != nil {
		return err
	}

	return s.waitForRollout(ctx)
}

func (s *StatelessOperator) Ingress(ctx context.Context) error { return s.base.ApplyIngress(ctx) }

func (k *StatelessOperator) Cleanup(ctx context.Context) {
	k.log.Debug("Running cleanup...")

	err := k.kubectl.Delete(ctx, "statefulset", k.config.AppIdentifier, k.config.AppNamespace)
	if err != nil && !errors.IsNotFound(err) {
		k.log.Debug("Cleanup: %v", err)
	}
}

// waitForRollout waits for the deployment to roll out using watch
func (s *StatelessOperator) waitForRollout(ctx context.Context) error {
	timeout := time.Duration(s.config.RolloutTimeout()) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return s.watchDeploymentRollout(ctx)
}

// watchDeploymentRollout watches a Deployment until it's ready
func (s *StatelessOperator) watchDeploymentRollout(ctx context.Context) error {
	// Get initial state
	deploy, err := s.clientset.AppsV1().Deployments(s.config.AppNamespace).Get(ctx, s.config.AppIdentifier, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	if isDeploymentReady(deploy) {
		return nil
	}

	// Watch for changes
	watcher, err := s.clientset.AppsV1().Deployments(s.config.AppNamespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", s.config.AppIdentifier),
	})
	if err != nil {
		return fmt.Errorf("failed to watch deployment: %w", err)
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Error("Rollout timed out, reverting...")
			s.rollbackDeployment(context.Background())
			return fmt.Errorf("deployment rollout timed out")
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}

			if event.Type == watch.Error {
				continue
			}

			deploy, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				continue
			}

			s.log.Debug("Deployment %s: %d/%d ready",
				deploy.Name,
				deploy.Status.ReadyReplicas,
				*deploy.Spec.Replicas)

			if isDeploymentReady(deploy) {
				return nil
			}

			// Check for failure conditions
			for _, cond := range deploy.Status.Conditions {
				if cond.Type == appsv1.DeploymentProgressing && cond.Status == corev1.ConditionFalse {
					s.log.Error("Deployment failed: %s", cond.Message)
					s.rollbackDeployment(context.Background())
					return fmt.Errorf("deployment failed: %s", cond.Message)
				}
			}
		}
	}
}

// rollbackDeployment rolls back a failed deployment using kubectl
func (k *StatelessOperator) rollbackDeployment(ctx context.Context) {
	args := []string{"rollout", "undo", "deployment", k.config.AppIdentifier, "-n", k.config.AppNamespace}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		k.log.Error("Failed to rollback: %s", strings.TrimSpace(string(output)))
	} else {
		k.log.Info("Rolled back deployment")
	}
}
