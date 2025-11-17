package convert

import (
	"fmt"

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
	step.AddScriptCmd(`echo "Git clone started"`)
	if in.Netrc != nil {
		step.AddScriptCmd(`echo "machine $GIT_MACHINE login $GIT_LOGIN password $GIT_PASSWORD" > ~/.netrc`)
		step.AddScriptCmd(`chmod 600 ~/.netrc`)
		addSecret(pCtx, step, "GIT_LOGIN", in.Netrc.Login)
		addSecret(pCtx, step, "GIT_PASSWORD", in.Netrc.Password)
		addSecret(pCtx, step, "GIT_MACHINE", in.Netrc.Machine)
	}
	step.AddScriptCmd(fmt.Sprintf(`git clone %s --branch %s %s;`, gitSource.RepoURL, gitSource.Branch, wsBuildVolumePath))
	if gitSource.Commit != "" {
		step.AddScriptCmd(fmt.Sprintf("git -C %s config advice.detachedHead false", wsBuildVolumePath))
		step.AddScriptCmd(fmt.Sprintf("git -C %s checkout %s", wsBuildVolumePath, gitSource.Commit))
	}
	step.AddScriptCmd(`echo "Git clone successful"`)

	//volume mounts
	step.VolumeMounts = append(step.VolumeMounts, getBuildVolumeMount(pCtx))

	return nil
}
