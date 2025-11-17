package variable

import (
	"context"

	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) Add(ctx context.Context, envID, appID int64, in *AddVariableInput) error {
	err := c.sanitizeCreateInput(in)
	if err != nil {
		return err
	}

	allVars, err := c.ListInEnvironment(ctx, envID)
	if err != nil {
		return err
	}

	newVar := newVariable(envID, appID, in.Key, in.Value, in.Type)
	if err := c.updateTextFromValue(ctx, envID, newVar, allVars); err != nil {
		return err
	}
	return c.upsert(ctx, newVar)
}

func (c *Controller) AddSystem(ctx context.Context, envID, appID int64, key, value string) error {
	sysVar := newVariable(envID, appID, key, value, enum.VariableTypeBuildAndRun)
	return c.upsert(ctx, sysVar)
}

func (c *Controller) sanitizeCreateInput(in *AddVariableInput) error {
	errors := check.NewValidationErrors()
	if err := check.VariableKey(in.Key); err != nil {
		errors.AddValidationError("new_key", err)
	}
	if IsSystemVar(in.Key) {
		errors.AddValidationError("new_value", check.NewValidationError("Cloudness platform variable cannot be overridden"))
	}
	if in.Value == "" {
		errors.AddValidationError("new_value", check.NewValidationError("Variable value is required"))
	}

	if errors.HasError() {
		return errors
	}
	return nil
}
