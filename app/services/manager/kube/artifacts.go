package kube

import (
	"bufio"
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *K8sManager) ListArtifacts(ctx context.Context, server *types.Server, app *types.Application) ([]*types.Artifact, error) {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return nil, err
	}

	var pods *corev1.PodList
	pods, err = client.CoreV1().Pods(app.Namespace()).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", app.GetIdentifierStr()),
	})

	log.Ctx(ctx).Debug().Any("pods", pods).Msg("Pods list")
	if err != nil {
		return nil, err
	}

	return m.ToArtifacts(app, pods), nil
}

func (m *K8sManager) TailLogs(ctx context.Context, server *types.Server, app *types.Application) (<-chan *types.ArtifactLogLine, <-chan error, error) {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return nil, nil, err
	}

	pods, err := client.CoreV1().Pods(app.Namespace()).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", app.GetIdentifierStr()),
	})
	if err != nil {
		return nil, nil, err
	}

	logc := make(chan *types.ArtifactLogLine)
	errc := make(chan error)

	go func() {
		for i, pod := range pods.Items {
			req := client.CoreV1().Pods(app.Namespace()).GetLogs(pod.Name, &corev1.PodLogOptions{
				Follow:    true,
				TailLines: nil,
			})
			podLogs, err := req.Stream(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("error streaming logs")
				errc <- err
				return
			}
			defer podLogs.Close()

			scanner := bufio.NewScanner(podLogs)
			for {
				select {
				case <-ctx.Done():
					log.Debug().Msg("closing artifacts log channel")
					return
				default:
					if scanner.Scan() {
						logc <- &types.ArtifactLogLine{
							ArtifactUID: fmt.Sprintf("%s-%d", app.Name, i),
							Log:         scanner.Text(),
						}
					}
				}
			}
		}
	}()

	return logc, errc, nil
}
