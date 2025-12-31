package convert

import (
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/pipeline"
	specSvc "github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func buildCommandNew(
	step *pipeline.Step,
	in *pipeline.RunnerContextInput,
	pCtx *pipeline.RunnerContext,
	spec *types.ApplicationSpec,
	buildVars map[string]string,
) error {
	if !specSvc.NeedsBuild(spec) {
		return nil
	}

	gitSource := spec.Build.Source.Git
	sourcePath := wsBuildVolumePath
	if gitSource.BasePath != "" && gitSource.BasePath != "/" {
		sourcePath = wsBuildVolumePath + gitSource.BasePath
	}

	image, _, cacheImage := specSvc.GetImage(in.Application, in.Deployment, in.Config)

	// Set common build environment variables
	addSecret(pCtx, step, "CLOUDNESS_BUILD_SOURCE_PATH", sourcePath)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_IMAGE", image)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_CACHE_IMAGE", cacheImage)
	addSecret(pCtx, step, "CLOUDNESS_IMAGE_REGISTRY", in.Config.PushRegistryURL)
	addSecret(pCtx, step, "CLOUDNESS_IMAGE_MIRROR_REGISTRY", in.Config.MirrorRegistryURL)

	if in.Config.MirrorRegistryEnabled && in.Config.MirrorRegistryURL != "" {
		addVariable(pCtx, step, "CLOUDNESS_MIRROR_ENABLED", "true")
	}

	// Set build args as space-separated key=value pairs
	if len(buildVars) > 0 {
		var args []string
		for k, v := range buildVars {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
		addSecret(pCtx, step, "CLOUDNESS_BUILD_ARGS", strings.Join(args, " "))
	}

	switch specSvc.GetBuilder(spec) {
	case enum.BuilderTypeDockerfile:
		addVariable(pCtx, step, "CLOUDNESS_BUILD_TYPE", "dockerfile")
		addVariable(pCtx, step, "CLOUDNESS_BUILD_DOCKERFILE", gitSource.Dockerfile)
	case enum.BuilderTypeNixpacks:
		addSecret(pCtx, step, "CLOUDNESS_BUILD_TYPE", "nixpacks")
		if gitSource.BuildCommand != "" {
			addVariable(pCtx, step, "CLOUDNESS_BUILD_CMD", gitSource.BuildCommand)
		}
		if spec.Deploy.StartCommand != "" {
			addVariable(pCtx, step, "CLOUDNESS_START_CMD", spec.Deploy.StartCommand)
		}
	}

	// Add image secrets for other steps
	addSecrets(pCtx, step, buildVars)

	// Run the build script
	step.AddScriptCmd(". /usr/local/lib/build-script.sh")

	return nil
}
