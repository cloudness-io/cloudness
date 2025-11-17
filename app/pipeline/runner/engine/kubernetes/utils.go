package kubernetes

import (
	"fmt"
	"os"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func isImagePullBackOffState(pod *v1.Pod) bool {
	for _, containerState := range pod.Status.ContainerStatuses {
		if isImagePullBackOffStateContainer(containerState) {
			return true
		}
	}
	return false
}

func isInvalidImageName(pod *v1.Pod) bool {
	for _, containerState := range pod.Status.ContainerStatuses {
		if isInvalidImageNameContainer(containerState) {
			return true
		}
	}
	return false
}

func isImagePullBackOffStateContainer(containerStatus v1.ContainerStatus) bool {
	return containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason == "ImagePullBackOff"
}

func isInvalidImageNameContainer(constainerStatus v1.ContainerStatus) bool {
	return constainerStatus.State.Waiting != nil && constainerStatus.State.Waiting.Reason == "InvalidImageName"
}

func parseMemoryGBtoQuantity(value float64) resource.Quantity {
	memValMiB := int64(value * 1024)

	memValStr := strconv.FormatInt(memValMiB, 10)
	return resource.MustParse(memValStr + "Mi")
}

func parseCPUtoQuantity(value float64) resource.Quantity {
	cpuValStr := fmt.Sprintf("%g", value)

	return resource.MustParse(cpuValStr)
}
