package variable

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

type UpdateVariableInput struct {
	Value string            `json:"value"`
	Type  enum.VariableType `json:"type"`
}

type GenerateVariableInput struct {
	Length int `json:"length,string"`
}

func (c *Controller) Update(ctx context.Context, envID, appID int64, varUID int64, in *UpdateVariableInput) error {
	variable, err := c.variableStore.Find(ctx, appID, varUID)
	if err != nil {
		return err
	}

	variable.Value = in.Value
	variable.Type = in.Type

	allVars, err := c.ListInEnvironment(ctx, envID)
	if err != nil {
		return err
	}
	if err := c.updateTextFromValue(ctx, envID, variable, allVars); err != nil {
		return err
	}

	return c.upsert(ctx, variable)
}

func (c *Controller) UpdateGenerate(ctx context.Context, envID, appID int64, varUID int64, in *GenerateVariableInput) error {
	variable, err := c.variableStore.Find(ctx, appID, varUID)
	if err != nil {
		return err
	}

	value, ref := c.generateSecret(in.Length)
	variable.TextValue = value
	variable.Value = ref
	return c.upsert(ctx, variable)
}

func (c *Controller) updateTextFromValue(ctx context.Context, envID int64, v *types.Variable, allVars []*types.Variable) error {
	ref := getRef(v.Value)
	//log.Ctx(ctx).Debug().Any("Ref", ref).Msg("Has reference")
	if len(ref) > 0 {
		nameToIDMap, idToValueMap := getNameAndIDRefMap(allVars, v.ApplicationID)
		v.Value = replaceWithMap(v.Value, nameToIDMap)
		v.TextValue = replaceWithMap(v.Value, idToValueMap)
		// log.Ctx(ctx).Debug().Any("Name_Ref", nameRef).Any("ID_Ref", idRef).Any("Value", v.Value).Any("Text", v.TextValue).Msg("All reference")
	} else {
		v.TextValue = v.Value
	}
	return nil
}
