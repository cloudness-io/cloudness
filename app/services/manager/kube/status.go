package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *K8sManager) ListApplicationStatuses(ctx context.Context, server *types.Server) ([]*types.AppStatus, error) {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	nowMilli := now.UnixMilli()
	statusByApp := make(map[int64]*types.AppStatus)

	deployments, err := client.AppsV1().Deployments(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=cloudness",
	})
	if err != nil {
		return nil, err
	}

	for _, deploy := range deployments.Items {
		appUID := m.getApplicationUIDFromPodLabels(deploy.Labels)
		projectID := m.getProjectIDFromPodLabels(deploy.Labels)
		updatedAt := m.getUpdateTimeFromPodAnnotations(deploy.Annotations)

		if appUID == 0 {
			continue
		}
		if updatedAt != 0 && updatedAt > (nowMilli-10_000) {
			continue
		}

		status, reason := evaluateDeploymentStatus(deploy)
		candidate := &types.AppStatus{
			Timestamp:      now,
			ApplicationUID: appUID,
			ProjectID:      projectID,
			InstanceName:   getApplicationNameFromLabels(deploy.Labels, deploy.Name),
			Status:         status,
			Reason:         reason,
		}

		mergeAppStatus(statusByApp, candidate)
	}

	statefulSets, err := client.AppsV1().StatefulSets(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=cloudness",
	})
	if err != nil {
		return nil, err
	}

	for _, sts := range statefulSets.Items {
		appUID := m.getApplicationUIDFromPodLabels(sts.Labels)
		projectID := m.getProjectIDFromPodLabels(sts.Labels)
		updatedAt := m.getUpdateTimeFromPodAnnotations(sts.Annotations)

		if appUID == 0 {
			continue
		}
		if updatedAt != 0 && updatedAt > (nowMilli-10_000) {
			continue
		}

		status, reason := evaluateStatefulSetStatus(sts)
		candidate := &types.AppStatus{
			Timestamp:      now,
			ApplicationUID: appUID,
			ProjectID:      projectID,
			InstanceName:   getApplicationNameFromLabels(sts.Labels, sts.Name),
			Status:         status,
			Reason:         reason,
		}

		mergeAppStatus(statusByApp, candidate)
	}

	pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=cloudness",
	})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		appUID := m.getApplicationUIDFromPodLabels(pod.Labels)
		projectID := m.getProjectIDFromPodLabels(pod.Labels)
		updatedAt := m.getUpdateTimeFromPodAnnotations(pod.Annotations)

		if appUID == 0 {
			continue
		}
		if updatedAt != 0 && updatedAt > (nowMilli-10_000) {
			continue
		}

		status, reason := evaluatePodStatus(pod)

		candidate := &types.AppStatus{
			Timestamp:      now,
			ApplicationUID: appUID,
			ProjectID:      projectID,
			InstanceName:   getApplicationNameFromLabels(pod.Labels, pod.Name),
			Status:         status,
			Reason:         reason,
		}

		mergeAppStatus(statusByApp, candidate)
	}

	statuses := make([]*types.AppStatus, 0, len(statusByApp))
	for _, s := range statusByApp {
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func evaluatePodStatus(pod corev1.Pod) (enum.ApplicationStatus, string) {
	// Check terminal phases first.
	switch pod.Status.Phase {
	case corev1.PodFailed:
		return enum.ApplicationStatusError, firstNonEmpty(pod.Status.Reason, pod.Status.Message, "pod failed")
	case corev1.PodUnknown:
		return enum.ApplicationStatusError, firstNonEmpty(pod.Status.Reason, pod.Status.Message, "pod unknown")
	case corev1.PodSucceeded:
		return enum.ApplicationStatusSleeping, "pod succeeded"
	}

	// For Pending and Running pods, inspect container-level statuses to detect
	// errors such as ImagePullBackOff, CrashLoopBackOff, etc. that the phase
	// alone does not surface.

	// Check init containers first — a failing init container blocks the whole pod.
	if status, reason, found := evaluateContainerStatuses(pod.Status.InitContainerStatuses); found {
		return status, reason
	}

	// Then check regular containers.
	if status, reason, found := evaluateContainerStatuses(pod.Status.ContainerStatuses); found {
		return status, reason
	}

	// No container-level error detected — fall back to phase.
	if pod.Status.Phase == corev1.PodPending {
		return enum.ApplicationStatusPaused, firstNonEmpty(pod.Status.Reason, pod.Status.Message, "pod pending")
	}

	return enum.ApplicationStatusRunning, ""
}

// evaluateContainerStatuses inspects a slice of container statuses and returns
// an error status if any container is in a bad state. The third return value
// indicates whether a conclusive status was found.
func evaluateContainerStatuses(statuses []corev1.ContainerStatus) (enum.ApplicationStatus, string, bool) {
	for _, cs := range statuses {
		if cs.State.Waiting != nil {
			reason := firstNonEmpty(cs.State.Waiting.Reason, cs.State.Waiting.Message, "container waiting")
			return enum.ApplicationStatusError, reason, true
		}

		if cs.State.Terminated != nil && cs.State.Terminated.ExitCode != 0 {
			reason := firstNonEmpty(cs.State.Terminated.Reason, cs.State.Terminated.Message, "container terminated")
			return enum.ApplicationStatusError, reason, true
		}

		if !cs.Ready && cs.State.Running != nil {
			return enum.ApplicationStatusError, "container not ready", true
		}
	}

	return "", "", false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func evaluateDeploymentStatus(deploy appsv1.Deployment) (enum.ApplicationStatus, string) {
	desired := int32(1)
	// if deploy.Spec.Replicas != nil {
	// 	desired = *deploy.Spec.Replicas
	// }

	// if desired == 0 {
	// 	return enum.ApplicationStatusSleeping, "scaled to zero"
	// }

	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentProgressing && cond.Status == corev1.ConditionFalse {
			return enum.ApplicationStatusError, firstNonEmpty(cond.Reason, cond.Message, "deployment not progressing")
		}
		if cond.Type == appsv1.DeploymentReplicaFailure && cond.Status == corev1.ConditionTrue {
			return enum.ApplicationStatusError, firstNonEmpty(cond.Reason, cond.Message, "deployment replica failure")
		}
	}

	if deploy.Status.ReadyReplicas < desired {
		return enum.ApplicationStatusPaused, fmt.Sprintf("ready replicas %d/%d", deploy.Status.ReadyReplicas, desired)
	}

	if deploy.Status.UpdatedReplicas < desired {
		return enum.ApplicationStatusPaused, fmt.Sprintf("updated replicas %d/%d", deploy.Status.UpdatedReplicas, desired)
	}

	return enum.ApplicationStatusRunning, ""
}

func evaluateStatefulSetStatus(sts appsv1.StatefulSet) (enum.ApplicationStatus, string) {
	desired := int32(1)
	// if sts.Spec.Replicas != nil {
	// 	desired = *sts.Spec.Replicas
	// }

	// if desired == 0 {
	// 	return enum.ApplicationStatusSleeping, "scaled to zero"
	// }

	for _, cond := range sts.Status.Conditions {
		if cond.Status == corev1.ConditionFalse {
			return enum.ApplicationStatusError, firstNonEmpty(cond.Reason, cond.Message, "statefulset condition failed")
		}
	}

	if sts.Status.ReadyReplicas < desired {
		return enum.ApplicationStatusPaused, fmt.Sprintf("ready replicas %d/%d", sts.Status.ReadyReplicas, desired)
	}

	if sts.Status.CurrentReplicas < desired {
		return enum.ApplicationStatusPaused, fmt.Sprintf("current replicas %d/%d", sts.Status.CurrentReplicas, desired)
	}

	return enum.ApplicationStatusRunning, ""
}

func mergeAppStatus(statusByApp map[int64]*types.AppStatus, candidate *types.AppStatus) {
	current, found := statusByApp[candidate.ApplicationUID]
	if !found || isHigherPriorityStatus(candidate.Status, current.Status) {
		statusByApp[candidate.ApplicationUID] = candidate
		return
	}

	if current.Status == candidate.Status && current.Reason == "" && candidate.Reason != "" {
		current.Reason = candidate.Reason
	}
}

func isHigherPriorityStatus(candidate enum.ApplicationStatus, current enum.ApplicationStatus) bool {
	priority := map[enum.ApplicationStatus]int{
		enum.ApplicationStatusError:    3,
		enum.ApplicationStatusPaused:   2,
		enum.ApplicationStatusRunning:  1,
		enum.ApplicationStatusSleeping: 0,
	}

	return priority[candidate] > priority[current]
}

func getApplicationNameFromLabels(labels map[string]string, fallback string) string {
	if name := labels["app.kubernetes.io/name"]; name != "" {
		return name
	}

	if name := labels["app.kubernetes.io/instance"]; name != "" {
		return name
	}

	if name := labels["app"]; name != "" {
		return name
	}

	return fallback
}
