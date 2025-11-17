package convert

import (
	"github.com/cloudness-io/cloudness/app/pipeline"
)

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
