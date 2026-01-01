package kubernetes

import (
	"fmt"

	"github.com/cloudness-io/cloudness/app/pipeline"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineIDLabel             = "pipeline.id"
	pipelineDeploymentUIDLabel  = "deployment.uid"
	pipelineApplicationUIDLabel = "application.uid"
	placeholderImage            = "cloudnessio/placeholder:2"
	pipelinePodLabel            = "runner.type"
	pipelinePodValue            = "kubernetes"
)

func toNameSpace(namespace string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
}

func toPodLabels(p *pipeline.RunnerContext) map[string]string {
	labels := make(map[string]string)

	labels[pipelineIDLabel] = p.RunnerName
	labels[pipelinePodLabel] = pipelinePodValue
	labels[pipelineDeploymentUIDLabel] = fmt.Sprintf("%d", p.Deployment.UID)
	labels[pipelineApplicationUIDLabel] = fmt.Sprintf("%d", p.ApplicationUID)

	return labels
}

func toSecret(p *pipeline.RunnerContext) *v1.Secret {
	stringData := make(map[string]string)
	for _, s := range p.Secrets {
		stringData[s.Name] = s.Data
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.RunnerName,
		},
		Type:       v1.SecretTypeOpaque,
		StringData: stringData,
	}
}

func toConfigMap(p *pipeline.RunnerContext) *v1.ConfigMap {
	stringData := make(map[string]string)
	for _, s := range p.Variables {
		stringData[s.Name] = s.Value
	}

	// Add ConfigFiles to ConfigMap
	for _, cf := range p.ConfigFiles {
		stringData[cf.Key] = cf.Content
	}

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.RunnerName,
		},
		Data: stringData,
	}
}

func toPod(nameSpace string, p *pipeline.RunnerContext) *v1.Pod {
	resource := toResources(p)
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.RunnerName,
			Namespace: nameSpace,
			Labels:    toPodLabels(p),
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "cloudness-runner-account",
			RestartPolicy:      v1.RestartPolicyNever,
			Volumes:            toVolume(p),
			Containers:         toContainers(p, resource),
			InitContainers:     toInitContainers(p),
			Resources:          &resource,
		},
	}
}

func toContainers(p *pipeline.RunnerContext, resource v1.ResourceRequirements) []v1.Container {
	var containers []v1.Container
	for _, step := range p.Steps {
		containers = append(containers, toContainer(step, p, resource))
	}

	return containers
}

func toContainer(s *pipeline.Step, p *pipeline.RunnerContext, resource v1.ResourceRequirements) v1.Container {
	return v1.Container{
		Name:            s.Name,
		Image:           placeholderImage,
		Command:         s.Command,
		Args:            s.Args,
		ImagePullPolicy: v1.PullIfNotPresent,
		WorkingDir:      s.WorkingDir,
		Env:             toEnv(s, p),
		VolumeMounts:    toVolumeMounts(s),
		ReadinessProbe:  toReadNessProbe(s),
		Resources:       resource,
		SecurityContext: &v1.SecurityContext{
			Privileged: &s.Privileged,
		},
	}
}

func toInitContainers(p *pipeline.RunnerContext) []v1.Container {
	var initContainers []v1.Container
	for _, step := range p.InitSteps {
		initContainers = append(initContainers, toInitContainer(step, p))
	}
	return initContainers
}

func toInitContainer(s *pipeline.Step, p *pipeline.RunnerContext) v1.Container {
	return v1.Container{
		Name:            s.Name,
		Image:           s.Image,
		Command:         s.Command,
		Args:            s.Args,
		ImagePullPolicy: v1.PullIfNotPresent,
		WorkingDir:      s.WorkingDir,
		Env:             toEnv(s, p),
		VolumeMounts:    toVolumeMounts(s),
	}
}

func toResources(p *pipeline.RunnerContext) v1.ResourceRequirements {
	var dst v1.ResourceRequirements
	if p.ResourcesLimit.CPU == 0 && p.ResourcesLimit.Memory == 0 {
		return dst
	}
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    parseCPUtoQuantity(p.ResourcesLimit.Memory),
			v1.ResourceMemory: parseMemoryGBtoQuantity(p.ResourcesLimit.CPU),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    parseCPUtoQuantity(0.1),
			v1.ResourceMemory: parseMemoryGBtoQuantity(0.1),
		},
	}
}

func toEnv(s *pipeline.Step, p *pipeline.RunnerContext) []v1.EnvVar {
	var envs []v1.EnvVar
	for k, v := range s.Envs {
		envs = append(envs, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	for _, secretEnv := range s.Secrets {
		envs = append(envs, v1.EnvVar{
			Name: secretEnv.Key,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: p.RunnerName,
					},
					Key: secretEnv.Key,
				},
			},
		})
	}

	for _, varEnv := range s.Variables {
		envs = append(envs, v1.EnvVar{
			Name: varEnv.Key,
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: p.RunnerName,
					},
					Key: varEnv.Key,
				},
			},
		})
	}
	return envs
}

func toVolume(p *pipeline.RunnerContext) []v1.Volume {
	var volumes []v1.Volume

	//workspace
	for _, workspace := range p.Workspaces {
		volumes = append(volumes, v1.Volume{
			Name: workspace.ID,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium: v1.StorageMediumMemory,
				},
			},
		})
	}

	// Add ConfigMap volume for file mounts if there are ConfigFiles
	if len(p.ConfigFiles) > 0 {
		volumes = append(volumes, v1.Volume{
			Name: "configfiles",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: p.RunnerName,
					},
				},
			},
		})
	}

	return volumes
}

func toVolumeMounts(s *pipeline.Step) []v1.VolumeMount {
	var volumeMounts []v1.VolumeMount
	for _, v := range s.VolumeMounts {
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      v.ID,
			MountPath: v.Path,
			ReadOnly:  v.Readonly,
		})
	}

	// Add ConfigFileMount volumes (mounts ConfigMap files to specified path)
	for _, cfm := range s.ConfigFileMounts {
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      "configfiles",
			MountPath: cfm.Path,
			ReadOnly:  true,
		})
	}

	return volumeMounts
}

func toReadNessProbe(s *pipeline.Step) *v1.Probe {
	var probe *v1.Probe
	if s.Liveness != nil {
		switch {
		case s.Liveness.Exec != nil:
			probe = &v1.Probe{
				ProbeHandler: v1.ProbeHandler{
					Exec: &v1.ExecAction{
						Command: s.Liveness.Exec.Command,
					},
				},
				InitialDelaySeconds: 5,
				PeriodSeconds:       5,
			}
		}
	}
	return probe
}
