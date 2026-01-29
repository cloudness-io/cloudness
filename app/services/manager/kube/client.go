package kube

import (
	"context"
	"os"

	"github.com/cloudness-io/cloudness/types"

	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// TODO: get server creds from server object
func (m *K8sManager) getInterface(ctx context.Context, server *types.Server) (kubernetes.Interface, error) {
	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient, err = getClientOutOfCluster()
	} else {
		kubeClient, err = getClientInCluster()
	}

	return kubeClient, err
}

func (m *K8sManager) getDynamicClient(ctx context.Context, server *types.Server) (*dynamic.DynamicClient, error) {
	config, err := m.getClientConfig(ctx, server)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func (m *K8sManager) getACMEClient(ctx context.Context, server *types.Server) (*cmclientset.Clientset, error) {
	config, err := m.getClientConfig(ctx, server)
	if err != nil {
		return nil, err
	}
	return cmclientset.NewForConfig(config)
}

func (m *K8sManager) getGatewayClient(ctx context.Context, server *types.Server) (*gatewayclientset.Clientset, error) {
	config, err := m.getClientConfig(ctx, server)
	if err != nil {
		return nil, err
	}
	return gatewayclientset.NewForConfig(config)
}

func (m *K8sManager) getMetricsClient(ctx context.Context, server *types.Server) (*metricsclient.Clientset, error) {
	config, err := m.getClientConfig(ctx, server)
	if err != nil {
		return nil, err
	}
	return metricsclient.NewForConfig(config)
}

func (m *K8sManager) getClientConfig(ctx context.Context, server *types.Server) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", getOutofClusterKubeConfigPath())
		if err != nil {
			return nil, err
		}
	}
	return config, err
}

func getOutofClusterKubeConfigPath() string {
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if len(kubeConfigPath) == 0 {
		kubeConfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	return kubeConfigPath
}

func getClientOutOfCluster() (kubernetes.Interface, error) {
	kubeConfigPath := getOutofClusterKubeConfigPath()

	//use current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getClientInCluster() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
