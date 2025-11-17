// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

// PVCExpansionRequest type info
var (
	pvrGVR = schema.GroupVersionResource{
		Group:    "cloudness.io",
		Version:  "v1alpha1",
		Resource: "pvcexpansionrequests",
	}
)

// Finalizer for cleanup logic
const (
	PVCExpansionFinalizer = "pvcexpansion.cloudness.io/finalizer"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// CRD is unstructured for brevity; not adding typed api to scheme.
}

type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// fetch unstructured PVCExpansionRequest
	var pvr unstructured.Unstructured
	pvr.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cloudness.io",
		Version: "v1alpha1",
		Kind:    "PVCExpansionRequest",
	})
	if err := r.Get(ctx, req.NamespacedName, &pvr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle deletion
	if pvr.GetDeletionTimestamp() != nil {
		return r.handleDeletion(ctx, &pvr)
	}

	// Add finalizer if not present
	if err := r.ensureFinalizer(ctx, &pvr); err != nil {
		log.Error(err, "failed to add finalizer")
		return ctrl.Result{}, err
	}

	// read spec fields
	spec, found, err := unstructured.NestedMap(pvr.Object, "spec")
	if err != nil || !found {
		_ = r.updateStatus(ctx, &pvr, "Failed", "spec not found or invalid")
		return ctrl.Result{}, fmt.Errorf("invalid spec")
	}

	// required fields: statefulSetRef.name and size
	ssRef, _, _ := unstructured.NestedMap(spec, "statefulSetRef")
	if ssRef == nil {
		_ = r.updateStatus(ctx, &pvr, "Failed", "statefulSetRef missing")
		return ctrl.Result{}, fmt.Errorf("statefulSetRef missing")
	}
	ssName, _, _ := unstructured.NestedString(ssRef, "name")
	ssNs, foundNs, _ := unstructured.NestedString(ssRef, "namespace")
	if !foundNs || ssNs == "" {
		ssNs = req.Namespace
	}
	if ssName == "" {
		_ = r.updateStatus(ctx, &pvr, "Failed", "statefulSetRef.name missing")
		return ctrl.Result{}, fmt.Errorf("statefulSetRef.name missing")
	}
	sizeStr, _, _ := unstructured.NestedString(spec, "size")
	if sizeStr == "" {
		_ = r.updateStatus(ctx, &pvr, "Failed", "size missing")
		return ctrl.Result{}, fmt.Errorf("size missing")
	}
	forceOffline, _, _ := unstructured.NestedBool(spec, "forceOffline")

	// Parse timeout from spec (default to 10m if not specified)
	timeoutStr, found, _ := unstructured.NestedString(spec, "timeout")
	if !found || timeoutStr == "" {
		timeoutStr = "10m"
	}
	timeout, err := r.parseTimeout(timeoutStr)
	if err != nil {
		_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("invalid timeout format: %v", err))
		return ctrl.Result{}, err
	}

	// update status: processing
	_ = r.updateStatus(ctx, &pvr, "Processing", fmt.Sprintf("Expanding %s/%s to %s (forceOffline=%v)", ssNs, ssName, sizeStr, forceOffline))

	// fetch statefulset
	var sts appsv1.StatefulSet
	if err := r.Get(ctx, client.ObjectKey{Namespace: ssNs, Name: ssName}, &sts); err != nil {
		_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("StatefulSet not found: %v", err))
		return ctrl.Result{}, err
	}

	origReplicas := int32(1)
	if sts.Spec.Replicas != nil {
		origReplicas = *sts.Spec.Replicas
	}

	// Build list of PVC names from volumeClaimTemplates
	// Note: assuming simple name pattern: <claimTemplateName>-<stsName>-<ordinal>
	claimTemplateNames := []string{}
	for _, t := range sts.Spec.VolumeClaimTemplates {
		claimTemplateNames = append(claimTemplateNames, t.ObjectMeta.Name)
	}
	if len(claimTemplateNames) == 0 {
		_ = r.updateStatus(ctx, &pvr, "Failed", "StatefulSet has no volumeClaimTemplates; controller currently supports only volumeClaimTemplates.")
		return ctrl.Result{}, fmt.Errorf("no volumeClaimTemplates")
	}

	// generate PVC names for each replica and claim template
	var pvcNames []types.NamespacedName
	for ordinal := int32(0); ordinal < origReplicas; ordinal++ {
		for _, ctn := range claimTemplateNames {
			// default naming convention
			pvcName := fmt.Sprintf("%s-%s-%d", ctn, sts.Name, ordinal)
			pvcNames = append(pvcNames, types.NamespacedName{Namespace: ssNs, Name: pvcName})
		}
	}

	// If forceOffline -> scale to 0
	if forceOffline && origReplicas > 0 {
		_ = r.updateStatus(ctx, &pvr, "ScalingDown", "Scaling StatefulSet to 0 for offline expansion")
		zero := int32(0)
		patch := client.MergeFrom(sts.DeepCopy())
		sts.Spec.Replicas = &zero
		if err := r.Patch(ctx, &sts, patch); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("failed to scale down: %v", err))
			return ctrl.Result{}, err
		}
		// wait for pods to terminate
		terminationTimeout := timeout / 3 // Use 1/3 of total timeout for pod termination
		if terminationTimeout < 2*time.Minute {
			terminationTimeout = 2 * time.Minute
		}
		if err := r.waitForPodsTermination(ctx, &sts, terminationTimeout); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("timed out waiting for pods to terminate: %v", err))
			return ctrl.Result{}, err
		}
	}

	// For each PVC, patch storage request
	_, err = resource.ParseQuantity(sizeStr)
	if err != nil {
		_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("invalid size: %v", err))
		return ctrl.Result{}, err
	}

	for _, pvcNN := range pvcNames {
		var pvc corev1.PersistentVolumeClaim
		if err := r.Get(ctx, pvcNN, &pvc); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("failed to get pvc %s/%s: %v", pvcNN.Namespace, pvcNN.Name, err))
			return ctrl.Result{}, err
		}

		// check storageClass allowVolumeExpansion
		if pvc.Spec.StorageClassName == nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("pvc %s has no storageClass - cannot determine expansion support", pvc.Name))
			return ctrl.Result{}, fmt.Errorf("pvc %s/%s has no storageClass", pvc.Namespace, pvc.Name)
		}
		scName := *pvc.Spec.StorageClassName

		// we will do a lightweight check by reading StorageClass object. Use dynamic since not typed imported here.
		// For brevity, we attempt to fetch storageclass via API server using Unstructured
		var scUnstructured unstructured.Unstructured
		scUnstructured.SetGroupVersionKind(schema.GroupVersionKind{Group: "storage.k8s.io", Version: "v1", Kind: "StorageClass"})
		if err := r.Get(ctx, types.NamespacedName{Name: scName}, &scUnstructured); err != nil {
			// Some clusters deny listing; if we cannot fetch, we will proceed but warn.
			ctrl.Log.Info("could not fetch StorageClass; proceeding", "storageClass", scName, "err", err)
		} else {
			allow, found, _ := unstructured.NestedBool(scUnstructured.Object, "allowVolumeExpansion")
			if !found || !allow {
				_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("storageClass %s does not allow expansion", scName))
				return ctrl.Result{}, fmt.Errorf("storageClass %s does not allow expansion", scName)
			}
		}

		// patch PVC
		patch := client.MergeFrom(pvc.DeepCopy())
		if pvc.Spec.Resources.Requests == nil {
			pvc.Spec.Resources.Requests = corev1.ResourceList{}
		}
		q, _ := resource.ParseQuantity(sizeStr)
		pvc.Spec.Resources.Requests[corev1.ResourceStorage] = q
		if err := r.Patch(ctx, &pvc, patch); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("failed to patch pvc %s: %v", pvc.Name, err))
			return ctrl.Result{}, err
		}
		_ = r.updateStatus(ctx, &pvr, "PatchingPVC", fmt.Sprintf("patched pvc %s to %s", pvc.Name, sizeStr))
	}

	// Wait for PV capacity or PVC conditions to reflect change
	_ = r.updateStatus(ctx, &pvr, "WaitingForResize", "Waiting for CSI/controller to complete expansion")
	resizeTimeout := timeout / 2 // Use half of total timeout for resize operations
	if resizeTimeout < 5*time.Minute {
		resizeTimeout = 5 * time.Minute
	}
	if err := r.waitForPVCsResize(ctx, pvcNames, sizeStr, resizeTimeout); err != nil {
		_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("resize did not complete: %v", err))
		return ctrl.Result{}, err
	}

	// Scale back up if we scaled down earlier
	if forceOffline && origReplicas > 0 {
		_ = r.updateStatus(ctx, &pvr, "ScalingUp", "Scaling StatefulSet back to original replica count")
		patch := client.MergeFrom(sts.DeepCopy())
		sts.Spec.Replicas = &origReplicas
		if err := r.Patch(ctx, &sts, patch); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("failed to scale up: %v", err))
			return ctrl.Result{}, err
		}
		// wait for pods to be ready
		readyTimeout := timeout / 3 // Use 1/3 of total timeout for pod readiness
		if readyTimeout < 5*time.Minute {
			readyTimeout = 5 * time.Minute
		}
		if err := r.waitForPodsReady(ctx, &sts, readyTimeout); err != nil {
			_ = r.updateStatus(ctx, &pvr, "Failed", fmt.Sprintf("timed out waiting for pods ready: %v", err))
			return ctrl.Result{}, err
		}
	}

	_ = r.updateStatus(ctx, &pvr, "Succeeded", "PVC expansion completed")

	// Remove finalizer on successful completion
	if err := r.removeFinalizer(ctx, &pvr); err != nil {
		log.Error(err, "failed to remove finalizer after successful completion")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// ensureFinalizer adds the finalizer if it's not already present
func (r *Reconciler) ensureFinalizer(ctx context.Context, pvr *unstructured.Unstructured) error {
	finalizers := pvr.GetFinalizers()
	for _, f := range finalizers {
		if f == PVCExpansionFinalizer {
			return nil // already present
		}
	}

	// Add finalizer
	finalizers = append(finalizers, PVCExpansionFinalizer)
	pvr.SetFinalizers(finalizers)

	return r.Update(ctx, pvr)
}

// removeFinalizer removes the finalizer from the resource
func (r *Reconciler) removeFinalizer(ctx context.Context, pvr *unstructured.Unstructured) error {
	finalizers := pvr.GetFinalizers()
	var newFinalizers []string

	for _, f := range finalizers {
		if f != PVCExpansionFinalizer {
			newFinalizers = append(newFinalizers, f)
		}
	}

	if len(newFinalizers) != len(finalizers) {
		pvr.SetFinalizers(newFinalizers)
		return r.Update(ctx, pvr)
	}

	return nil
}

// handleDeletion handles cleanup when the resource is being deleted
func (r *Reconciler) handleDeletion(ctx context.Context, pvr *unstructured.Unstructured) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("handling deletion of PVCExpansionRequest")

	// Check if our finalizer is present
	hasFinalizer := false
	for _, f := range pvr.GetFinalizers() {
		if f == PVCExpansionFinalizer {
			hasFinalizer = true
			break
		}
	}

	if !hasFinalizer {
		return ctrl.Result{}, nil // nothing to do
	}

	// Get current status to understand what cleanup is needed
	phase, _, _ := unstructured.NestedString(pvr.Object, "status", "phase")

	// Perform cleanup based on current phase
	if err := r.performCleanup(ctx, pvr, phase); err != nil {
		log.Error(err, "cleanup failed, will retry")
		_ = r.updateStatus(ctx, pvr, "Failed", fmt.Sprintf("Cleanup failed: %v", err))
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Update status to indicate cleanup completion
	_ = r.updateStatus(ctx, pvr, "Succeeded", "Cleanup completed successfully")

	// Remove our finalizer
	if err := r.removeFinalizer(ctx, pvr); err != nil {
		log.Error(err, "failed to remove finalizer during deletion")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	log.Info("successfully completed cleanup and removed finalizer")
	return ctrl.Result{}, nil
}

// performCleanup performs the actual cleanup operations
func (r *Reconciler) performCleanup(ctx context.Context, pvr *unstructured.Unstructured, phase string) error {
	log := ctrl.LoggerFrom(ctx)

	// Read spec to get StatefulSet reference
	spec, found, err := unstructured.NestedMap(pvr.Object, "spec")
	if err != nil || !found {
		return fmt.Errorf("invalid spec during cleanup")
	}

	ssRef, _, _ := unstructured.NestedMap(spec, "statefulSetRef")
	if ssRef == nil {
		return nil // nothing to cleanup
	}

	ssName, _, _ := unstructured.NestedString(ssRef, "name")
	ssNs, foundNs, _ := unstructured.NestedString(ssRef, "namespace")
	if !foundNs || ssNs == "" {
		ssNs = pvr.GetNamespace()
	}

	if ssName == "" {
		return nil // nothing to cleanup
	}

	// Fetch the StatefulSet to check its current state
	var sts appsv1.StatefulSet
	if err := r.Get(ctx, client.ObjectKey{Namespace: ssNs, Name: ssName}, &sts); err != nil {
		if errors.IsNotFound(err) {
			log.Info("StatefulSet not found during cleanup, assuming already cleaned up")
			return nil
		}
		return fmt.Errorf("failed to fetch StatefulSet during cleanup: %w", err)
	}

	// Cleanup logic based on the phase where deletion was requested
	switch phase {
	case "ScalingDown":
		// If we were scaling down, try to restore original replica count
		return r.restoreOriginalReplicas(ctx, &sts, pvr)

	case "ScalingUp":
		// If we were scaling up, the StatefulSet should already be at correct scale
		log.Info("cleanup during ScalingUp phase, StatefulSet should be at correct scale")
		return nil

	case "PatchingPVC", "WaitingForResize":
		// PVC patches are atomic and don't need rollback
		// The expansion either succeeded or failed, but we don't rollback storage
		log.Info("cleanup during PVC operations, no rollback needed for storage changes")
		return nil

	case "Processing", "Pending":
		// Early phase, likely no changes made yet
		log.Info("cleanup during early phase, no significant changes to rollback")
		return nil

	default:
		log.Info("cleanup for unknown phase, performing general cleanup", "phase", phase)
		return nil
	}
}

// restoreOriginalReplicas attempts to restore the StatefulSet to its original replica count
func (r *Reconciler) restoreOriginalReplicas(ctx context.Context, sts *appsv1.StatefulSet, pvr *unstructured.Unstructured) error {
	log := ctrl.LoggerFrom(ctx)

	// Try to get original replica count from status
	expandedPVCs, found, _ := unstructured.NestedSlice(pvr.Object, "status", "expandedPVCs")
	if !found || len(expandedPVCs) == 0 {
		log.Info("no original replica information found, leaving StatefulSet as-is")
		return nil
	}

	// If StatefulSet is already scaled to 0, try to restore it
	currentReplicas := int32(0)
	if sts.Spec.Replicas != nil {
		currentReplicas = *sts.Spec.Replicas
	}

	if currentReplicas == 0 {
		// Calculate original replicas from the number of PVCs we were processing
		// This is a best-effort restoration
		originalReplicas := int32(len(expandedPVCs) / len(sts.Spec.VolumeClaimTemplates))
		if originalReplicas > 0 {
			log.Info("restoring StatefulSet replica count", "from", currentReplicas, "to", originalReplicas)

			patch := client.MergeFrom(sts.DeepCopy())
			sts.Spec.Replicas = &originalReplicas

			if err := r.Patch(ctx, sts, patch); err != nil {
				return fmt.Errorf("failed to restore StatefulSet replicas: %w", err)
			}
		}
	}

	return nil
}

func (r *Reconciler) updateStatus(ctx context.Context, u *unstructured.Unstructured, phase, message string) error {
	if err := unstructured.SetNestedField(u.Object, phase, "status", "phase"); err != nil {
		return err
	}
	if err := unstructured.SetNestedField(u.Object, message, "status", "message"); err != nil {
		return err
	}
	// set observedGeneration
	if err := unstructured.SetNestedField(u.Object, int64(u.GetGeneration()), "status", "observedGeneration"); err != nil {
		return err
	}
	return r.Status().Update(ctx, u)
}

func (r *Reconciler) waitForPodsTermination(ctx context.Context, sts *appsv1.StatefulSet, timeout time.Duration) error {
	// wait until no pods for this sts
	start := time.Now()
	for {
		var podList corev1.PodList
		selector := fields.SelectorFromSet(map[string]string{"statefulset.kubernetes.io/pod-name": ""}) // dummy
		_ = selector                                                                                    // not used; do labelSelector based on sts.Spec.Selector
		if err := r.List(ctx, &podList, client.InNamespace(sts.Namespace), client.MatchingLabels(sts.Spec.Selector.MatchLabels)); err != nil {
			return err
		}
		if len(podList.Items) == 0 {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for pods termination; %d pods still present", len(podList.Items))
		}
		time.Sleep(2 * time.Second)
	}
}

func (r *Reconciler) waitForPodsReady(ctx context.Context, sts *appsv1.StatefulSet, timeout time.Duration) error {
	start := time.Now()
	for {
		var podList corev1.PodList
		if err := r.List(ctx, &podList, client.InNamespace(sts.Namespace), client.MatchingLabels(sts.Spec.Selector.MatchLabels)); err != nil {
			return err
		}
		ready := 0
		for _, p := range podList.Items {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
					ready++
				}
			}
		}
		if int32(ready) == *sts.Spec.Replicas {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for pods ready; ready=%d target=%d", ready, *sts.Spec.Replicas)
		}
		time.Sleep(3 * time.Second)
	}
}

func (r *Reconciler) waitForPVCsResize(ctx context.Context, pvcNNs []types.NamespacedName, target string, timeout time.Duration) error {
	start := time.Now()
	targetQ, _ := resource.ParseQuantity(target)
	for {
		allGood := true
		for _, nn := range pvcNNs {
			var pvc corev1.PersistentVolumeClaim
			if err := r.Get(ctx, nn, &pvc); err != nil {
				return err
			}
			// check bound PV capacity via pvc.Status.Capacity
			capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]
			if ok && capacity.Cmp(targetQ) >= 0 {
				continue
			}
			allGood = false
			break
		}
		if allGood {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for pvc capacity to be >= %s", target)
		}
		time.Sleep(3 * time.Second)
	}
}

// parseTimeout parses timeout string like "10m", "1h", "30s" into time.Duration
func (r *Reconciler) parseTimeout(timeoutStr string) (time.Duration, error) {
	if timeoutStr == "" {
		return 10 * time.Minute, nil // default
	}

	// Parse the timeout string
	if len(timeoutStr) < 2 {
		return 0, fmt.Errorf("invalid timeout format: %s", timeoutStr)
	}

	unit := timeoutStr[len(timeoutStr)-1:]
	valueStr := timeoutStr[:len(timeoutStr)-1]

	// Parse the numeric value
	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return 0, fmt.Errorf("invalid timeout value: %s", valueStr)
	}

	// Convert based on unit
	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid timeout unit: %s (use s, m, or h)", unit)
	}
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election")
	flag.Parse()

	ctrl.SetLogger(ctrl.Log.WithName("pvc-expansion-controller"))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	reconciler := &Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}

	if err = ctrl.NewControllerManagedBy(mgr).
		For(&unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "cloudness.io/v1alpha1",
			"kind":       "PVCExpansionRequest",
		}}).
		Complete(reconciler); err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
