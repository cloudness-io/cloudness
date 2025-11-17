package kubernetes

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
)

const (
	defaultResyncDuration = time.Second * 5
)

type kube struct {
	nameSpace string
	client    kubernetes.Interface
	config    types.RunnerConfig
}

func New() engine.Engine {
	return &kube{}
}

func (e *kube) Type() string {
	return "kubernetes"
}

func (e *kube) IsAvailable(_ types.RunnerConfig) bool {
	if len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0 {
		return true
	}
	kubeConfigPath := getOutofClusterKubeConfigPath()

	_, err := os.Stat(kubeConfigPath)
	return err == nil
}

func (e *kube) Load(ctx context.Context, config types.RunnerConfig) (*engine.EngineInfo, error) {
	e.config = config
	e.nameSpace = config.KubeNameSpace

	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient, err = getClientOutOfCluster()
	} else {
		kubeClient, err = getClientInCluster()
	}

	if err != nil {
		return nil, err
	}

	e.client = kubeClient

	_, err = e.client.CoreV1().Namespaces().Create(ctx, toNameSpace(e.nameSpace), metav1.CreateOptions{})
	if err != nil && !kerror.IsAlreadyExists(err) {
		return nil, err
	}

	return &engine.EngineInfo{}, nil
}

func (e *kube) ListIncomplete(ctx context.Context) ([]int64, error) {
	pods, err := e.client.CoreV1().Pods(e.nameSpace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", pipelinePodLabel, pipelinePodValue),
	})
	if err != nil {
		return nil, err
	}

	dst := make([]int64, len(pods.Items))
	for i, p := range pods.Items {
		uidStr := p.Labels[pipelineDeploymentUIDLabel]
		if uid, err := strconv.ParseInt(uidStr, 10, 64); err != nil {
			log.Ctx(ctx).Error().Err(err).Str(pipelineDeploymentUIDLabel, uidStr).Msg("engine: unable up parse deployment uid in resume phase")
		} else {
			dst[i] = uid
		}
	}

	return dst, nil
}

func (e *kube) Setup(rCtx context.Context, pCtx *pipeline.RunnerContext) error {
	_, err := e.client.CoreV1().Secrets(e.nameSpace).Create(rCtx, toSecret(pCtx), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = e.client.CoreV1().Pods(e.nameSpace).Create(rCtx, toPod(e.nameSpace, pCtx), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return e.tailPod(rCtx, pCtx.RunnerName)
}

func (e *kube) StartStep(rCtx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) error {
	log := zerolog.Ctx(rCtx)
	log.Trace().Msg("engine: starting step")

	var backoff = wait.Backoff{
		Steps:    15,
		Duration: 500 * time.Millisecond,
		Factor:   1.0,
		Jitter:   0.5,
	}

	//flag to determine if the step container found
	found := false

	err := retry.RetryOnConflict(backoff, func() error {
		//TODO: Lock mutex?
		pod, err := e.client.CoreV1().Pods(e.nameSpace).Get(rCtx, pCtx.RunnerName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		for i, container := range pod.Spec.Containers {
			if container.Name != step.Name {
				//TODO: what if container not found?
				continue
			}

			found = true

			pod.Spec.Containers[i].Image = step.Image

			pod.Labels["pipeline.step"] = step.Name
			_, err = e.client.CoreV1().Pods(e.nameSpace).Update(rCtx, pod, metav1.UpdateOptions{})
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("engine: container with name %s not found", step.Name)
	}

	return e.tailPod(rCtx, pCtx.RunnerName)
}

func (e *kube) TailStep(ctx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) (io.ReadCloser, error) {
	if err := e.tailContainer(ctx, pCtx.RunnerName, step.Name); err != nil {
		return nil, err
	}

	opts := &v1.PodLogOptions{
		Follow:    true,
		Container: step.Name,
	}

	logs, err := e.client.CoreV1().RESTClient().Get().
		Namespace(e.nameSpace).
		Name(pCtx.RunnerName).
		Resource("pods").
		SubResource("log").
		VersionedParams(opts, scheme.ParameterCodec).
		Stream(ctx)
	if err != nil {
		return nil, err
	}

	rc, wc := io.Pipe()

	go func() {
		defer logs.Close()
		defer wc.Close()
		defer rc.Close()

		_, err = io.Copy(wc, logs)
		if err != nil {
			return
		}
	}()

	return rc, nil
}

func (e *kube) WaitStep(ctx context.Context, pCtx *pipeline.RunnerContext, step *pipeline.Step) (*engine.State, error) {
	if err := e.tailPod(ctx, pCtx.RunnerName); err != nil {
		return nil, err
	}

	pod, err := e.client.CoreV1().Pods(e.nameSpace).Get(ctx, pCtx.RunnerName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(pod.Status.ContainerStatuses) == 0 {
		return nil, fmt.Errorf("no container status found in pod %s", pCtx.RunnerName)
	}

	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name != step.Name {
			continue
		}

		if cs.State.Terminated == nil {
			return nil, fmt.Errorf("no terminated state found in container %s/%s", pCtx.RunnerName, step.Name)
		}
		return &engine.State{
			ExitCode:  int(cs.State.Terminated.ExitCode),
			Exited:    true,
			OOMKilled: false,
		}, nil
	}

	return nil, fmt.Errorf("no status found for container %s/%s", pCtx.RunnerName, step.Name)
}

func (e *kube) Destroy(ctx context.Context, pCtx *pipeline.RunnerContext) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Any("Panic", r).Msg("engine: error destroying engine artifacts")
		}
	}()

	//delete secret
	err := e.client.CoreV1().Secrets(e.nameSpace).Delete(ctx, pCtx.RunnerName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return e.client.CoreV1().Pods(e.nameSpace).Delete(ctx, pCtx.RunnerName, metav1.DeleteOptions{})
}

func (e *kube) tailPod(rCtx context.Context, podName string) error {
	log := zerolog.Ctx(rCtx)

	var backoff = wait.Backoff{
		Steps:    10,
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0.1,
		Cap:      30 * time.Second,
	}

	return retry.RetryOnConflict(backoff, func() error {
		upChan := make(chan bool)
		errChan := make(chan error)
		timeoutChan := make(chan bool)

		// Set a timeout for this retry attempt
		timeout := time.NewTimer(30 * time.Second)
		defer timeout.Stop()
		go func() {
			<-timeout.C
			timeoutChan <- true
		}()

		podUpdate := func(_, new any) {
			pod, ok := new.(*v1.Pod)
			if !ok {
				return
			}

			if pod.Name == podName {
				if isImagePullBackOffState(pod) || isInvalidImageName(pod) {
					errChan <- fmt.Errorf("could not pull image for deployment %s", podName)
					return
				}
				switch pod.Status.Phase {
				case v1.PodRunning, v1.PodSucceeded:
					upChan <- true
				case v1.PodFailed:
					errChan <- fmt.Errorf("pod failed: %s", podName)
				}
			}
		}

		labelOptions := informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = fmt.Sprintf("%s=%s", pipelineIDLabel, podName)
		})
		si := informers.NewSharedInformerFactoryWithOptions(e.client, defaultResyncDuration, informers.WithNamespace(e.nameSpace), labelOptions)
		_, err := si.Core().V1().Pods().Informer().AddEventHandler(
			cache.ResourceEventHandlerFuncs{
				UpdateFunc: podUpdate,
			},
		)
		if err != nil {
			log.Warn().Err(err).Msg("failed to add event handler, retrying...")
			return err
		}

		stopper := make(chan struct{})
		defer close(stopper)
		go si.Start(stopper)

		// Wait for informer cache to sync
		if !cache.WaitForCacheSync(rCtx.Done(), si.Core().V1().Pods().Informer().HasSynced) {
			return fmt.Errorf("failed to sync informer cache")
		}

		for {
			select {
			case <-rCtx.Done():
				return rCtx.Err()
			case <-upChan:
				return nil
			case err := <-errChan:
				// Don't retry for permanent failures like image pull errors or pod failures
				if err != nil {
					log.Error().Err(err).Msg("permanent failure in tailPod")
					return wait.ErrorInterrupted(err)
				}
				return err
			case <-timeoutChan:
				log.Warn().Str("podName", podName).Msg("timeout waiting for pod, retrying...")
				return fmt.Errorf("timeout waiting for pod %s", podName)
			}
		}
	})
}

func (e *kube) tailContainer(ctx context.Context, podName string, containerName string) error {
	log := log.Ctx(ctx)
	upChan := make(chan bool)
	errChan := make(chan error)

	podUpdate := func(_, new any) {
		pod, ok := new.(*v1.Pod)
		if !ok {
			log.Trace().Msg("could not parse pod while trailing")
			return
		}

		if pod.Name == podName {
			if pod.Status.Phase == v1.PodPending {
				return
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == containerName {
					if isImagePullBackOffStateContainer(containerStatus) || isInvalidImageNameContainer(containerStatus) {
						errChan <- fmt.Errorf("could not pull image for container %s", containerName)
					}
					if containerStatus.State.Waiting != nil {
						return
					}
					upChan <- true
				}
			}
		}
	}

	labelOptions := informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
		lo.LabelSelector = fmt.Sprintf("%s=%s", pipelineIDLabel, podName)
	})
	si := informers.NewSharedInformerFactoryWithOptions(e.client, defaultResyncDuration, informers.WithNamespace(e.nameSpace), labelOptions)
	_, err := si.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: podUpdate,
		},
	)
	if err != nil {
		return err
	}

	stopper := make(chan struct{})
	defer close(stopper)
	go si.Start(stopper)

	for {
		select {
		case <-upChan:
			return nil
		case err = <-errChan:
			return err
		}
	}
}
