package convert

import (
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/services/config"
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

	switch specSvc.GetBuilder(spec) {
	case enum.BuilderTypeDockerfile:
		updateDockerFileBuildStep(in, pCtx, spec, step, buildVars)
	case enum.BuilderTypeNixpacks:
		updateNixpacksBuildStep(in, pCtx, spec, step, buildVars)
	}

	addImageSecrets(in, pCtx, step)
	addSecrets(pCtx, step, buildVars)

	return nil
}

func updateDockerFileBuildStep(
	in *pipeline.RunnerContextInput,
	pCtx *pipeline.RunnerContext,
	spec *types.ApplicationSpec,
	step *pipeline.Step,
	buildVars map[string]string,
) {
	var cmd strings.Builder
	gitSource := spec.Build.Source.Git
	sourcePath := wsBuildVolumePath
	if gitSource.BasePath != "" && gitSource.BasePath != "/" {
		sourcePath = wsBuildVolumePath + gitSource.BasePath
	}
	image, _, cacheImage := specSvc.GetImage(in.Application, in.Deployment, in.Config)

	addBuildKitConfig(&cmd, in.Config)

	// Construct the buildctl-daemonless.sh command
	cmd.WriteString(fmt.Sprintf(`buildctl-daemonless.sh build \
        --frontend=dockerfile.v0 \
        --local context=%[1]s \
        --local dockerfile=%[1]s \
		  --opt filename=%[2]s \
        --output type=image,name=%[3]s,push=true `, sourcePath, gitSource.Dockerfile, image))

	if cacheImage != "" {
		cmd.WriteString(fmt.Sprintf(`--export-cache type=registry,ref=%s,mode=max `, cacheImage))
		cmd.WriteString(fmt.Sprintf(`--import-cache type=registry,ref=%s,mode=max `, cacheImage))
	}

	// Add BuildArgs
	if len(buildVars) > 0 {
		for k, v := range buildVars {
			cmd.WriteString(fmt.Sprintf(`--opt build-arg:%s=%s `, k, v)) //
		}
	}

	step.AddScriptCmd(cmd.String())
}

func updateNixpacksBuildStep(
	in *pipeline.RunnerContextInput,
	pCtx *pipeline.RunnerContext,
	spec *types.ApplicationSpec,
	step *pipeline.Step,
	buildVars map[string]string,
) {
	gitSource := spec.Build.Source.Git
	sourcePath := wsBuildVolumePath
	if gitSource.BasePath != "" && gitSource.BasePath != "/" {
		sourcePath = wsBuildVolumePath + gitSource.BasePath
	}

	image, _, cacheImage := specSvc.GetImage(in.Application, in.Deployment, in.Config)
	nixCommand := []string{fmt.Sprintf("nixpacks build %[1]s -o %[1]s", sourcePath)}
	nixCommand = append(nixCommand, fmt.Sprintf("--name %s", image))
	if gitSource.BuildCommand != "" {
		nixCommand = append(nixCommand, fmt.Sprintf(`--build-cmd "%s"`, gitSource.BuildCommand))
	}
	if spec.Deploy.StartCommand != "" {
		nixCommand = append(nixCommand, fmt.Sprintf(`--start-cmd "%s"`, spec.Deploy.StartCommand))
	}

	for key, value := range buildVars {
		nixCommand = append(nixCommand, fmt.Sprintf(`--env %s="%s"`, key, value))
	}
	nixCommand = append(nixCommand, "--verbose")
	step.AddScriptCmd(strings.Join(nixCommand, ` `))

	var cmd strings.Builder
	//generate buildkit toml
	addBuildKitConfig(&cmd, in.Config)

	// Construct the buildctl-daemonless.sh command
	cmd.WriteString(fmt.Sprintf(`buildctl-daemonless.sh build \
        --frontend=dockerfile.v0 \
        --local context=%[1]s \
        --local dockerfile=%[1]s \
		  --opt filename=/.nixpacks/Dockerfile \
        --output type=image,name=%[2]s,push=true `, sourcePath, image))

	if cacheImage != "" {
		cmd.WriteString(fmt.Sprintf(`--export-cache type=registry,ref=%s `, cacheImage))
		cmd.WriteString(fmt.Sprintf(`--import-cache type=registry,ref=%s,mode=max `, cacheImage))
	}

	step.AddScriptCmd(cmd.String())

}

func addImageSecrets(in *pipeline.RunnerContextInput, pCtx *pipeline.RunnerContext, step *pipeline.Step) {
	image, pullImage, cacheImage := specSvc.GetImage(in.Application, in.Deployment, in.Config)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_IMAGE", image)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_PULL_IMAGE", pullImage)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_CACHE_IMAGE", cacheImage)
	addSecret(pCtx, step, "CLOUDNESS_IMAGE_REGISTRY", in.Config.PushRegistryURL)
	addSecret(pCtx, step, "CLOUDNESS_IMAGE_MIRROR_REGISTRY", in.Config.MirrorRegistryURL)
}

func addBuildKitConfig(cmd *strings.Builder, config *config.PipelineConfig) {
	//generate buildkit toml
	cmd.WriteString(`BUILDKITD_CONFIG_PATH="$HOME/.config/buildkit/buildkitd.toml"` + "\n")
	cmd.WriteString(`mkdir -p "$(dirname "$BUILDKITD_CONFIG_PATH")"` + "\n")
	cmd.WriteString(`> "$BUILDKITD_CONFIG_PATH"` + "\n")
	cmd.WriteString(fmt.Sprintf(`MAIN_REGISTRY=$(echo %s | cut -d'/' -f1)`+"\n", config.PushRegistryURL))
	cmd.WriteString(`echo "[registry.\"$MAIN_REGISTRY\"]" >> "$BUILDKITD_CONFIG_PATH"` + "\n")
	cmd.WriteString(`echo "  http = true" >> "$BUILDKITD_CONFIG_PATH"` + "\n")
	cmd.WriteString(`echo "  insecure = true" >> "$BUILDKITD_CONFIG_PATH"` + "\n")

	if config.MirrorRegistryEnabled && config.MirrorRegistryURL != "" {
		// Add the mirror configuration
		cmd.WriteString(fmt.Sprintf(`MIRROR_REGISTRY=$(echo %s | cut -d'/' -f1)`+"\n", config.MirrorRegistryURL))
		cmd.WriteString(`echo "[registry.\"$MIRROR_REGISTRY\"]" >> "$BUILDKITD_CONFIG_PATH"` + "\n")
		cmd.WriteString(`echo "  http = true" >> "$BUILDKITD_CONFIG_PATH"` + "\n")
		cmd.WriteString(`echo "  insecure = true" >> "$BUILDKITD_CONFIG_PATH"` + "\n")

		cmd.WriteString(`echo "[registry.\"docker.io\"]" >> "$BUILDKITD_CONFIG_PATH"` + "\n")
		cmd.WriteString(fmt.Sprintf(`echo "  mirrors = [\"%s\"]" >> "$BUILDKITD_CONFIG_PATH"`+"\n", config.MirrorRegistryURL))
	}
}
