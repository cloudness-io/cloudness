package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type StatefulOperator struct {
	config    *Config
	base      *BaseOpeator
	clientset kubernetes.Interface
	kubectl   *Kubectl
	log       *Logger
}

func NewStatefulOperator(cfg *Config, clientset kubernetes.Interface, log *Logger) *StatefulOperator {
	kubectl := NewKubectl(cfg, log)
	baseOperator := NewBaseOperator(cfg, kubectl, log)
	return &StatefulOperator{
		config:    cfg,
		base:      baseOperator,
		clientset: clientset,
		kubectl:   kubectl,
		log:       log,
	}
}

func (s *StatefulOperator) ApplyCommon(ctx context.Context) error {
	return s.base.ApplyCommon(ctx)
}

func (s *StatefulOperator) Volumes(ctx context.Context) error {
	if s.config.HasVolume {
		if err := s.deployVolumes(ctx); err != nil {
			return fmt.Errorf("failed to deploy volumes: %w", err)
		}
		s.log.Step("Volumes provisioned")
	}
	return nil
}

func (k *StatefulOperator) Deploy(ctx context.Context) error {
	if err := k.kubectl.ApplyYAMLFile(ctx, k.config.AppYAMLPath); err != nil {
		return err
	}

	return k.waitForRollout(ctx)
}

func (s *StatefulOperator) Ingress(ctx context.Context) error { return s.base.ApplyIngress(ctx) }

func (k *StatefulOperator) Cleanup(ctx context.Context) {
	k.log.Debug("Running cleanup...")

	err := k.kubectl.Delete(ctx, "deployment", k.config.AppIdentifier, k.config.AppNamespace)
	if err != nil && !errors.IsNotFound(err) {
		k.log.Debug("Cleanup: %v", err)
	}
}

// waitForRollout waits for the statefulset to roll out using watch
func (k *StatefulOperator) waitForRollout(ctx context.Context) error {
	timeout := time.Duration(k.config.RolloutTimeout()) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return k.watchStatefulSetRollout(ctx)
}

// watchStatefulSetRollout watches a StatefulSet until it's ready
func (k *StatefulOperator) watchStatefulSetRollout(ctx context.Context) error {
	// Get initial state
	sts, err := k.clientset.AppsV1().StatefulSets(k.config.AppNamespace).Get(ctx, k.config.AppIdentifier, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get statefulset: %w", err)
	}

	if isStatefulSetReady(sts) {
		return nil
	}

	// Watch for changes
	watcher, err := k.clientset.AppsV1().StatefulSets(k.config.AppNamespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", k.config.AppIdentifier),
	})
	if err != nil {
		return fmt.Errorf("failed to watch statefulset: %w", err)
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			k.log.Error("Rollout timed out")
			return fmt.Errorf("statefulset rollout timed out")
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}

			if event.Type == watch.Error {
				continue
			}

			sts, ok := event.Object.(*appsv1.StatefulSet)
			if !ok {
				continue
			}

			k.log.Debug("StatefulSet %s: %d/%d ready",
				sts.Name,
				sts.Status.ReadyReplicas,
				*sts.Spec.Replicas)

			if isStatefulSetReady(sts) {
				return nil
			}
		}
	}
}

// =============================================================================
// Volume Handling (client-go for PVC status checking)
// =============================================================================

// deployVolumes handles PVC creation and resizing
func (k *StatefulOperator) deployVolumes(ctx context.Context) error {
	// Handle remount for volume resize
	if k.config.NeedRemount {
		k.log.Info("Volume resize detected, removing statefulset for remount...")
		_ = k.kubectl.Delete(ctx, "statefulset", k.config.AppIdentifier, k.config.AppNamespace)
	}

	// Apply volume configuration
	if err := k.kubectl.ApplyYAMLFile(ctx, k.config.VolumeYAMLPath); err != nil {
		return err
	}

	// Wait for PVCs using watch (efficient)
	return k.waitForPVCs(ctx)
}

// waitForPVCs waits for all PVCs to be bound or resized
func (k *StatefulOperator) waitForPVCs(ctx context.Context) error {
	// List PVCs with our label
	pvcs, err := k.clientset.CoreV1().PersistentVolumeClaims(k.config.AppNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", k.config.AppIdentifier),
	})
	if err != nil {
		return fmt.Errorf("failed to list PVCs: %w", err)
	}

	if len(pvcs.Items) == 0 {
		return nil
	}

	timeout := time.Duration(k.config.PVCResizeTimeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, pvc := range pvcs.Items {
		if err := k.waitForPVC(ctx, pvc.Name); err != nil {
			return err
		}
	}

	return nil
}

// waitForPVC waits for a single PVC to be ready using watch
func (k *StatefulOperator) waitForPVC(ctx context.Context, name string) error {
	// First check current status
	pvc, err := k.clientset.CoreV1().PersistentVolumeClaims(k.config.AppNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get PVC %s: %w", name, err)
	}

	if k.isPVCReady(pvc) {
		return nil
	}

	// Watch for changes
	watcher, err := k.clientset.CoreV1().PersistentVolumeClaims(k.config.AppNamespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	})
	if err != nil {
		return fmt.Errorf("failed to watch PVC %s: %w", name, err)
	}
	defer watcher.Stop()

	k.log.Info("Waiting for volume '%s'...", name)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for PVC %s", name)
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed for PVC %s", name)
			}

			if event.Type == watch.Error {
				continue
			}

			pvc, ok := event.Object.(*corev1.PersistentVolumeClaim)
			if !ok {
				continue
			}

			k.log.Debug("PVC %s status: %s", name, pvc.Status.Phase)

			if k.isPVCReady(pvc) {
				return nil
			}
		}
	}
}

// isPVCReady checks if a PVC is ready for use
func (k *StatefulOperator) isPVCReady(pvc *corev1.PersistentVolumeClaim) bool {
	// Bound is ready
	if pvc.Status.Phase == corev1.ClaimBound {
		return true
	}

	// WaitForFirstConsumer - consider ready (will bind when pod schedules)
	if pvc.Status.Phase == corev1.ClaimPending {
		if pvc.Spec.StorageClassName != nil {
			// Check for WaitForFirstConsumer annotation
			if ann := pvc.Annotations; ann != nil {
				if ann["volume.kubernetes.io/selected-node"] != "" {
					return true
				}
			}
		}
	}

	// FileSystemResizePending - volume resized, waiting for pod remount
	for _, cond := range pvc.Status.Conditions {
		if cond.Type == corev1.PersistentVolumeClaimFileSystemResizePending && cond.Status == corev1.ConditionTrue {
			k.log.Info("Volume %s resized. Remounting application to finalize.", pvc.Name)
			return true
		}
	}

	return false
}
