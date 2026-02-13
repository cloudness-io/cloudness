package kube

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *K8sManager) ListMetrics(ctx context.Context, server *types.Server) ([]*types.AppMetrics, error) {
	metricsClient, err := m.getMetricsClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=cloudness",
	})
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	metrics := make([]*types.AppMetrics, 0)

	for _, pod := range podMetrics.Items {
		appUID := m.getApplicationUIDFromPodLabels(pod.Labels)
		if appUID == 0 {
			continue
		}
		var cpu int64
		var mem int64

		for _, container := range pod.Containers {
			cpu += container.Usage.Cpu().MilliValue()
			mem += container.Usage.Memory().Value()
		}

		metrics = append(metrics, &types.AppMetrics{
			Timestamp:      now,
			ApplicationUID: appUID,
			InstanceName:   pod.Name,
			CPU:            cpu,
			Memory:         mem,
		})
	}

	return metrics, nil
}
