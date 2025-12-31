package convert

import (
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/types"
)

func getRunnerNamespace(in *pipeline.RunnerContext) string {
	return fmt.Sprintf("cloudness-runner-%d-%d",
		in.ApplicationUID,
		in.Deployment.UID)
}

func getBuildWorksapceVolumeId(in *pipeline.RunnerContext) string {
	return fmt.Sprintf("%s-build-workspace", getRunnerNamespace(in))
}

func getDeploymentWorkspaceVolumeId(in *pipeline.RunnerContext) string {
	return fmt.Sprintf("%s-deployment-workspace", getRunnerNamespace(in))
}

func getBuildVolumeMount(in *pipeline.RunnerContext) *pipeline.VolumeMount {
	return &pipeline.VolumeMount{
		ID:       getBuildWorksapceVolumeId(in),
		Path:     wsBuildVolumePath,
		Readonly: false,
	}
}

func getDeployVolumeMount(in *pipeline.RunnerContext) *pipeline.VolumeMount {
	return &pipeline.VolumeMount{
		ID:       getDeploymentWorkspaceVolumeId(in),
		Path:     wsDeployVolumePath,
		Readonly: false,
	}
}

// var replacer
func replaceEnvVars(input string, vars map[string]*types.Variable) string {
	var builder strings.Builder
	for i := 0; i < len(input); {
		if input[i] == '$' && i+1 < len(input) && input[i+1] == '{' {
			end := strings.IndexByte(input[i:], '}')
			if end != -1 {
				key := input[i+2 : i+end]
				if val, ok := vars[key]; ok {
					builder.WriteString(val.TextValue)
				} else {
					builder.WriteString("${" + key + "}")
				}
				i += end + 1
				continue
			}
		}
		builder.WriteByte(input[i])
		i++
	}

	return builder.String()
}

func addSecrets(ctx *pipeline.RunnerContext, step *pipeline.Step, vars map[string]string) {
	for key, val := range vars {
		addSecret(ctx, step, key, val)
	}
}

func addSecret(ctx *pipeline.RunnerContext, step *pipeline.Step, key string, value string) {
	ctx.Secrets = append(ctx.Secrets, &pipeline.Secret{
		Name: key,
		Data: value,
		Mask: true,
	})

	step.Secrets = append(step.Secrets, &pipeline.SecretEnv{
		Key: key,
	})
}

func addVariable(ctx *pipeline.RunnerContext, step *pipeline.Step, key string, value string) {
	ctx.Variables = append(ctx.Variables, &pipeline.Variable{
		Name:  key,
		Value: value,
	})

	step.Variables = append(step.Variables, &pipeline.VariableEnv{
		Key: key,
	})
}
