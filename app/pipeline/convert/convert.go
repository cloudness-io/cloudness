package convert

import (
	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/types/enum"
)

func ToRunnerContext(in *pipeline.RunnerContextInput) (*pipeline.RunnerContext, error) {
	pCtx := &pipeline.RunnerContext{
		Deployment:     in.Deployment,
		ApplicationUID: in.Application.UID,
		Domain:         in.Application.Domain,
		Steps:          []*pipeline.Step{},
		InitSteps:      []*pipeline.Step{},
		Secrets:        []*pipeline.Secret{},
	}

	pCtx.RunnerName = getRunnerNamespace(pCtx)
	pCtx.Workspaces = []*pipeline.Volume{
		{
			ID: getBuildWorksapceVolumeId(pCtx),
		},
		{
			ID: getDeploymentWorkspaceVolumeId(pCtx),
		},
	}

	spec := pCtx.Deployment.Spec
	// build and run command
	if spec.HasBuildCommand() {
		spec.Build.Source.Git.BuildCommand = replaceEnvVars(spec.Build.Source.Git.BuildCommand, in.Variables)
	}
	if spec.HasStartCommand() {
		spec.Deploy.StartCommand = replaceEnvVars(spec.Deploy.StartCommand, in.Variables)
	}

	// args
	buildVars := map[string]string{}
	runVars := map[string]string{}
	for _, v := range in.Variables {
		switch v.Type {
		case enum.VariableTypeBuild:
			buildVars[v.Key] = v.TextValue
		case enum.VariableTypeRun:
			runVars[v.Key] = v.TextValue
		case enum.VariableTypeBuildAndRun:
			buildVars[v.Key] = v.TextValue
			runVars[v.Key] = v.TextValue
		}
	}

	step := &pipeline.Step{
		Name:           "deploy",
		Image:          baseImage,
		Command:        []string{"/bin/sh", "-c"},
		ScriptCommands: []string{},
		VolumeMounts:   []*pipeline.VolumeMount{},
		Privileged:     true,
		Envs:           map[string]string{},
	}

	step.AddStripCmds("#!/bin/sh\n\n", "set -e")

	if err := initCommand(step, in, pCtx, spec); err != nil {
		return nil, err
	}

	if err := buildCommandNew(step, in, pCtx, spec, buildVars); err != nil {
		return nil, err
	}

	if err := deployCommand(step, in, pCtx, spec, runVars); err != nil {
		return nil, err
	}
	step.Args = []string{step.GenerateShellScript()}

	pCtx.Steps = append(pCtx.Steps, step)

	return pCtx, nil
}
