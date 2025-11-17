package kube

import (
	"context"

	"github.com/cloudness-io/cloudness/types"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *K8sManager) GetIP(ctx context.Context, server *types.Server) (string, error) {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return "", err
	}

	svc, err := client.CoreV1().Services(defaultK8sGatewayNamespace).Get(ctx, defaultK8sGatewayService, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if svc.Spec.Type == v1.ServiceTypeLoadBalancer && len(svc.Status.LoadBalancer.Ingress) > 0 {
		return svc.Status.LoadBalancer.Ingress[0].IP, nil
	}

	return "", nil
}
