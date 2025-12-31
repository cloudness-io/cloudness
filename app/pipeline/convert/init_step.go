package convert

import (
	"github.com/cloudness-io/cloudness/app/pipeline"
	specSvc "github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/types"
)

func initCommand(
	step *pipeline.Step,
	in *pipeline.RunnerContextInput,
	pCtx *pipeline.RunnerContext,
	spec *types.ApplicationSpec,
) error {
	if !specSvc.NeedsInit(spec) {
		return nil
	}

	gitSource := spec.Build.Source.Git

	// Set environment variables for the init script
	addVariable(pCtx, step, "CLOUDNESS_GIT_REPO_URL", gitSource.RepoURL)
	addVariable(pCtx, step, "CLOUDNESS_GIT_BRANCH", gitSource.Branch)
	addVariable(pCtx, step, "CLOUDNESS_GIT_COMMIT", gitSource.Commit)
	addSecret(pCtx, step, "CLOUDNESS_BUILD_PATH", wsBuildVolumePath)

	// Set netrc credentials if provided
	if in.Netrc != nil {
		addSecret(pCtx, step, "GIT_LOGIN", in.Netrc.Login)
		addSecret(pCtx, step, "GIT_PASSWORD", in.Netrc.Password)
		addSecret(pCtx, step, "GIT_MACHINE", in.Netrc.Machine)
	}

	// Run the init script
	step.AddScriptCmd(". /usr/local/lib/init-script.sh")

	// Volume mounts
	step.VolumeMounts = append(step.VolumeMounts, getBuildVolumeMount(pCtx))

	return nil
}
