package convert

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/pipeline/convert/templates"
	specSvc "github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	shlex "github.com/kballard/go-shellquote"
)

func deployCommand(
	step *pipeline.Step,
	in *pipeline.RunnerContextInput,
	pCtx *pipeline.RunnerContext,
	spec *types.ApplicationSpec,
	vars map[string]string,
) error {
	namespace := in.Application.Namespace()
	name := in.Application.GetIdentifierStr()
	_, pullImage, _ := specSvc.GetImage(in.Application, in.Deployment, in.Config)

	tmplIn, err := getTemplateInput(pullImage, in, spec, vars)
	if err != nil {
		return err
	}

	common, volume, app, route, err := templates.GenerateKubeTemplates(tmplIn)
	if err != nil {
		return err
	}

	// Add YAML files as ConfigFiles (will be mounted from ConfigMap)
	addConfigFile(pCtx, "common.yaml", common)
	addConfigFile(pCtx, "volume.yaml", volume)
	addConfigFile(pCtx, "app.yaml", app)
	addConfigFile(pCtx, "route.yaml", route)

	// Mount ConfigFiles to deploy path
	step.ConfigFileMounts = append(step.ConfigFileMounts, &pipeline.ConfigFileMount{
		Path: wsDeployVolumePath,
		Keys: []string{"common.yaml", "volume.yaml", "app.yaml", "route.yaml"},
	})

	// General secrets for log sanitization
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_APP_IDENTIFIER", name)
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_APP_NAMESPACE", namespace)

	// Flags
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_FLAG_APP_TYPE", string(in.Application.Type))
	if len(tmplIn.Volumes) > 0 {
		addSecret(pCtx, step, "CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME", "1")
	}
	if tmplIn.ServiceDomain != nil {
		addSecret(pCtx, step, "CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE", "1")
	}

	// Deploy path for script to find YAML files
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_PATH", wsDeployVolumePath)

	// Common deployment info
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_TARGET_NAMESPACE", namespace)
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_TARGET_NAME", name)
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_TARGET_IMAGE", pullImage)

	// Unmount before update volumes?
	needsRemount := "0"
	if in.ServerRestctions.UnmountBeforeResize {
		if needsVolumeRemount(in.Deployment, in.PreviousDeployment) {
			needsRemount = "1"
		}
	}
	addSecret(pCtx, step, "CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT", needsRemount)

	// Run the Go-based deployer binary
	step.AddScriptCmd("cloudness-deploy")

	return nil
}

func getTemplateInput(image string, input *pipeline.RunnerContextInput, spec *types.ApplicationSpec, vars map[string]string) (*templates.TemplateIn, error) {
	in := &templates.TemplateIn{
		Namespace:   input.Application.Namespace(),
		Identifier:  input.Application.GetIdentifierStr(),
		Image:       image,
		MaxReplicas: spec.Deploy.MaxReplicas,
		CPU:         spec.Deploy.CPU,
		Memory:      spec.Deploy.Memory,
		Volumes:     make([]*templates.Volume, 0),
		Variables:   make(map[string]string),
		Secrets:     make(map[string]string),
		//networking
		ServicePorts: make([]int, 0),
		HasState:     input.Application.Type == enum.ApplicationTypeStateful,
	}

	if spec.ShouldAddStartCommand() {
		parts, err := shlex.Split(spec.Deploy.StartCommand)
		if err != nil {
			return nil, err
		}

		if len(parts) > 0 {
			in.Command = parts[:1]
			in.Args = parts[1:]
		}
	}

	for key, value := range vars {
		//NOTE: hack to avoid deployment failure due to empty value in secrets
		if value != "" {
			in.Secrets[key] = base64.StdEncoding.EncodeToString([]byte(value))
		} else {
			in.Variables[key] = `""`
		}
	}

	containerPorts := make(map[int]bool, 0)
	//allcontainer ports
	for _, port := range spec.Networking.ContainerPorts {
		containerPorts[port] = true
	}

	// service domain
	if spec.Networking.ServiceDomain != nil {
		containerPorts[spec.Networking.ServiceDomain.Port] = true
		url, err := url.Parse(spec.Networking.ServiceDomain.Domain)
		if err != nil {
			return nil, err
		}
		in.ServiceDomain = &templates.ServiceDomain{
			Websecure: url.Scheme == "https",
			Domain:    url.Hostname(),
			Port:      spec.Networking.ServiceDomain.Port,
		}
	}

	if spec.Networking.TCPProxies != nil {
		containerPorts[spec.Networking.TCPProxies.TCPPort] = true
	}

	for port := range containerPorts {
		in.ServicePorts = append(in.ServicePorts, port)
	}

	for _, v := range input.Volumes {
		in.Volumes = append(in.Volumes, &templates.Volume{
			VolumeName: v.GetIdentifierStr(),
			Storage:    fmt.Sprintf("%dGi", v.Size),
			MountPath:  v.MountPath,
		})
	}

	return in, nil
}

func needsVolumeRemount(currDeployment *types.Deployment, prevDeployment *types.Deployment) bool {
	if prevDeployment == nil {
		return false
	}
	currSpec := currDeployment.Spec
	prevSpec := prevDeployment.Spec
	if prevSpec == nil || len(prevSpec.Volumes) == 0 || len(currSpec.Volumes) == 0 {
		return false
	}

	if len(currSpec.Volumes) != len(prevSpec.Volumes) {
		return true
	}

	for i := 0; i < len(currSpec.Volumes); i++ {
		if currSpec.Volumes[i].VolumeSize != prevSpec.Volumes[i].VolumeSize {
			return true
		}
	}

	return false
}

// addConfigFile adds a file to the ConfigMap for mounting
func addConfigFile(pCtx *pipeline.RunnerContext, filename string, content string) {
	if content == "" {
		return
	}
	pCtx.ConfigFiles = append(pCtx.ConfigFiles, &pipeline.ConfigFile{
		Key:      filename,
		Filename: filename,
		Content:  content,
	})
}
