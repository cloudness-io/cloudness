package kube

import (
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	corev1 "k8s.io/api/core/v1"
)

func (m *K8sManager) ToArtifacts(app *types.Application, pods *corev1.PodList) []*types.Artifact {
	artifacts := make([]*types.Artifact, len(pods.Items))
	for i, pod := range pods.Items {
		artifacts[i] = &types.Artifact{
			UID:  pod.Name,
			Name: fmt.Sprintf("%s-%d", app.Name, i+1),
		}
	}
	return artifacts
}
