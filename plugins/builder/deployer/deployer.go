package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Deployer handles Kubernetes deployments
type Deployer struct {
	config        *Config
	clientset     kubernetes.Interface
	dynamicClient dynamic.Interface
	mapper        meta.RESTMapper
	log           *Logger
}

// NewDeployer creates a new deployer instance
func NewDeployer(cfg *Config) (*Deployer, error) {
	// Get in-cluster config
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Create discovery client and mapper
	dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return &Deployer{
		config:        cfg,
		clientset:     clientset,
		dynamicClient: dynamicClient,
		mapper:        mapper,
		log:           NewLogger(cfg.Verbose),
	}, nil
}

// Deploy runs the full deployment workflow
func (d *Deployer) Deploy(ctx context.Context) error {
	d.log.Section("Deploying application")

	// Step 1: Apply common artifacts (namespace, service account, etc.)
	if err := d.applyCommonArtifacts(ctx); err != nil {
		return fmt.Errorf("failed to apply common artifacts: %w", err)
	}
	d.log.Step("Prerequisite artifacts configured")

	// Step 2: Handle volumes (if any)
	if d.config.HasVolume {
		if err := d.deployVolumes(ctx); err != nil {
			return fmt.Errorf("failed to deploy volumes: %w", err)
		}
		d.log.Step("Volumes provisioned")
	}

	// Step 3: Deploy application
	if err := d.deployApplication(ctx); err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}
	d.log.Step("Application deployed")

	// Step 4: Configure routes (if any)
	if d.config.HasRoute {
		if err := d.deployRoutes(ctx); err != nil {
			return fmt.Errorf("failed to deploy routes: %w", err)
		}
		d.log.Step("HTTP routes configured")
	}

	// Cleanup: Remove opposite resource type
	d.cleanup(ctx)

	d.log.Success("Deployment completed successfully!")
	return nil
}

// applyCommonArtifacts applies common Kubernetes resources
func (d *Deployer) applyCommonArtifacts(ctx context.Context) error {
	return d.applyYAMLFile(ctx, d.config.CommonYAMLPath)
}

// deployVolumes handles PVC creation and resizing
func (d *Deployer) deployVolumes(ctx context.Context) error {
	// Handle remount for volume resize
	if d.config.NeedRemount {
		d.log.Info("Volume resize detected, removing statefulset for remount...")
		if err := d.deleteStatefulSet(ctx); err != nil {
			d.log.Warn("Failed to remove statefulset: %v", err)
		}
	}

	// Apply volume configuration
	if err := d.applyYAMLFile(ctx, d.config.VolumeYAMLPath); err != nil {
		return err
	}

	// Wait for PVCs to be ready
	return d.waitForPVCs(ctx)
}

// deployApplication deploys the main application and waits for rollout
func (d *Deployer) deployApplication(ctx context.Context) error {
	if err := d.applyYAMLFile(ctx, d.config.AppYAMLPath); err != nil {
		return err
	}

	return d.waitForRollout(ctx)
}

// deployRoutes applies route/ingress configuration
func (d *Deployer) deployRoutes(ctx context.Context) error {
	return d.applyYAMLFile(ctx, d.config.RouteYAMLPath)
}

// applyYAMLFile reads and applies a YAML file containing one or more manifests
func (d *Deployer) applyYAMLFile(ctx context.Context, filepath string) error {
	// Check if file exists and has content
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) || (err == nil && info.Size() == 0) {
		d.log.Debug("Skipping empty or non-existent file: %s", filepath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filepath, err)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	return d.applyYAML(ctx, data)
}

// applyYAML applies YAML content containing one or more Kubernetes manifests
func (d *Deployer) applyYAML(ctx context.Context, yamlData []byte) error {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bufio.NewReader(bytes.NewReader(yamlData)), 4096)
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode yaml: %w", err)
		}

		if rawObj.Raw == nil {
			continue
		}

		obj := &unstructured.Unstructured{}
		_, gvk, err := decUnstructured.Decode(rawObj.Raw, nil, obj)
		if err != nil {
			return fmt.Errorf("failed to decode object: %w", err)
		}

		// Get the REST mapping for this GVK
		mapping, err := d.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("failed to get REST mapping for %v: %w", gvk, err)
		}

		// Get the resource interface
		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			ns := obj.GetNamespace()
			if ns == "" {
				ns = d.config.AppNamespace
				obj.SetNamespace(ns)
			}
			dr = d.dynamicClient.Resource(mapping.Resource).Namespace(ns)
		} else {
			dr = d.dynamicClient.Resource(mapping.Resource)
		}

		// Apply using server-side apply
		d.log.Debug("Applying %s/%s in %s", gvk.Kind, obj.GetName(), obj.GetNamespace())
		_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, rawObj.Raw, metav1.PatchOptions{
			FieldManager: "cloudness-deployer",
		})
		if err != nil {
			// Fall back to Create/Update
			_, err = dr.Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = dr.Update(ctx, obj, metav1.UpdateOptions{})
			}
			if err != nil {
				return fmt.Errorf("failed to apply %s/%s: %w", gvk.Kind, obj.GetName(), err)
			}
		}
	}

	return nil
}

// waitForRollout waits for the deployment/statefulset to roll out
func (d *Deployer) waitForRollout(ctx context.Context) error {
	timeout := time.Duration(d.config.RolloutTimeout()) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.log.Error("Rollout timed out after %v, reverting...", timeout)
			d.rollback(context.Background()) // Use background context for rollback
			return fmt.Errorf("rollout timed out")
		case <-ticker.C:
			ready, err := d.checkRolloutStatus(ctx)
			if err != nil {
				d.log.Debug("Error checking rollout status: %v", err)
				continue
			}
			if ready {
				return nil
			}
		}
	}
}

// checkRolloutStatus checks if the deployment/statefulset is ready
func (d *Deployer) checkRolloutStatus(ctx context.Context) (bool, error) {
	if d.config.AppType == AppTypeStateless {
		deploy, err := d.clientset.AppsV1().Deployments(d.config.AppNamespace).Get(ctx, d.config.AppIdentifier, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return d.isDeploymentReady(deploy), nil
	}

	sts, err := d.clientset.AppsV1().StatefulSets(d.config.AppNamespace).Get(ctx, d.config.AppIdentifier, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return d.isStatefulSetReady(sts), nil
}

func (d *Deployer) isDeploymentReady(deploy *appsv1.Deployment) bool {
	if deploy.Generation != deploy.Status.ObservedGeneration {
		return false
	}
	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			if deploy.Status.UpdatedReplicas == *deploy.Spec.Replicas &&
				deploy.Status.ReadyReplicas == *deploy.Spec.Replicas {
				return true
			}
		}
	}
	return false
}

func (d *Deployer) isStatefulSetReady(sts *appsv1.StatefulSet) bool {
	if sts.Generation != sts.Status.ObservedGeneration {
		return false
	}
	return sts.Status.UpdatedReplicas == *sts.Spec.Replicas &&
		sts.Status.ReadyReplicas == *sts.Spec.Replicas
}

// rollback rolls back a failed deployment
func (d *Deployer) rollback(ctx context.Context) {
	if d.config.AppType == AppTypeStateless {
		// For deployments, we can use rollout undo
		deploys := d.clientset.AppsV1().Deployments(d.config.AppNamespace)
		deploy, err := deploys.Get(ctx, d.config.AppIdentifier, metav1.GetOptions{})
		if err != nil {
			d.log.Error("Failed to get deployment for rollback: %v", err)
			return
		}

		// Set revision to previous
		if deploy.Annotations == nil {
			deploy.Annotations = make(map[string]string)
		}
		deploy.Annotations["deployment.kubernetes.io/revision"] = ""

		_, err = deploys.Update(ctx, deploy, metav1.UpdateOptions{})
		if err != nil {
			d.log.Error("Failed to rollback deployment: %v", err)
		}
	}
	// StatefulSets don't support automatic rollback
}

// waitForPVCs waits for all PVCs to be bound or resized
func (d *Deployer) waitForPVCs(ctx context.Context) error {
	// Read PVC names from the volume YAML
	data, err := os.ReadFile(d.config.VolumeYAMLPath)
	if err != nil {
		return nil // No volume file, skip
	}

	pvcNames := d.extractPVCNames(data)
	if len(pvcNames) == 0 {
		return nil
	}

	timeout := time.Duration(d.config.PVCResizeTimeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, pvcName := range pvcNames {
		if err := d.waitForPVC(ctx, pvcName); err != nil {
			return err
		}
	}

	return nil
}

// extractPVCNames extracts PVC names from YAML content
func (d *Deployer) extractPVCNames(yamlData []byte) []string {
	var names []string
	decoder := yamlutil.NewYAMLOrJSONDecoder(bufio.NewReader(bytes.NewReader(yamlData)), 4096)

	for {
		var obj unstructured.Unstructured
		if err := decoder.Decode(&obj); err != nil {
			break
		}
		if obj.GetKind() == "PersistentVolumeClaim" {
			names = append(names, obj.GetName())
		}
	}
	return names
}

// waitForPVC waits for a single PVC to be ready
func (d *Deployer) waitForPVC(ctx context.Context, name string) error {
	ticker := time.NewTicker(time.Duration(d.config.PVCResizePollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for PVC %s", name)
		case <-ticker.C:
			pvc, err := d.clientset.CoreV1().PersistentVolumeClaims(d.config.AppNamespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				d.log.Debug("Error getting PVC %s: %v", name, err)
				continue
			}

			d.log.Debug("PVC %s status: %s", name, pvc.Status.Phase)

			switch pvc.Status.Phase {
			case corev1.ClaimBound:
				return nil
			case corev1.ClaimPending:
				// Check for WaitForFirstConsumer
				if d.isWaitingForFirstConsumer(ctx, name) {
					d.log.Debug("PVC %s waiting for first consumer", name)
					return nil
				}
			}

			// Check for FileSystemResizePending
			for _, cond := range pvc.Status.Conditions {
				if cond.Type == corev1.PersistentVolumeClaimFileSystemResizePending && cond.Status == corev1.ConditionTrue {
					d.log.Info("Volume %s resized. Remounting application to finalize.", name)
					return nil
				}
			}
		}
	}
}

// isWaitingForFirstConsumer checks if PVC is pending due to WaitForFirstConsumer
func (d *Deployer) isWaitingForFirstConsumer(ctx context.Context, pvcName string) bool {
	events, err := d.clientset.CoreV1().Events(d.config.AppNamespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.kind=PersistentVolumeClaim,involvedObject.name=%s", pvcName),
	})
	if err != nil {
		return false
	}

	for _, event := range events.Items {
		if event.Reason == "WaitForFirstConsumer" {
			return true
		}
	}
	return false
}

// deleteStatefulSet deletes the StatefulSet (for volume remount)
func (d *Deployer) deleteStatefulSet(ctx context.Context) error {
	return d.clientset.AppsV1().StatefulSets(d.config.AppNamespace).Delete(
		ctx,
		d.config.AppIdentifier,
		metav1.DeleteOptions{},
	)
}

// cleanup removes the opposite resource type (cleanup old deployment/statefulset)
func (d *Deployer) cleanup(ctx context.Context) {
	d.log.Debug("Running cleanup...")

	var err error
	if d.config.AppType == AppTypeStateless {
		// Remove any StatefulSet with the same name
		err = d.clientset.AppsV1().StatefulSets(d.config.AppNamespace).Delete(
			ctx, d.config.AppIdentifier, metav1.DeleteOptions{},
		)
	} else {
		// Remove any Deployment with the same name
		err = d.clientset.AppsV1().Deployments(d.config.AppNamespace).Delete(
			ctx, d.config.AppIdentifier, metav1.DeleteOptions{},
		)
	}

	if err != nil && !errors.IsNotFound(err) && !strings.Contains(err.Error(), "not found") {
		d.log.Warn("Failed to cleanup opposite resource: %v", err)
	}
}
