package variable

import (
	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ToDTO(allVars []*types.Variable, appID int64) *types.VariableDTO {
	sysVars := make([]*types.SystemVariable, 0)
	usrVars := make([]*types.UserVariable, 0)

	idRef := getIDToNameMap(allVars, appID)
	for _, v := range allVars {
		if v.ApplicationID == appID { //is from application in the current context
			switch true {
			case IsSystemVar(v.Key):
				sysVars = append(sysVars, &types.SystemVariable{
					Key:       v.Key,
					Type:      string(v.Type),
					TextValue: v.TextValue,
				})
			case isSecret(v.Value):
				usrVars = append(usrVars, &types.UserVariable{
					UID:            v.UID,
					Key:            v.Key,
					Type:           string(v.Type),
					TextValue:      v.TextValue,
					ReferenceValue: v.TextValue,
					IsGenerated:    true,
				})
			default:
				usrVars = append(usrVars, &types.UserVariable{
					UID:            v.UID,
					Key:            v.Key,
					Type:           string(v.Type),
					TextValue:      v.TextValue,
					ReferenceValue: replaceWithMap(v.Value, idRef),
				})
			}
		}
	}
	return &types.VariableDTO{
		SystemVariable: sysVars,
		UserVariable:   usrVars,
	}
}

// toListUpdate returns list of variables with updated text value for the application in context.
// Also returns the list of variables that need to be updated
func (c *Controller) toListUpdate(allVars []*types.Variable, appID int64) (map[string]*types.Variable, []*types.Variable) {
	updatedVars := make([]*types.Variable, 0)
	appVars := make(map[string]*types.Variable)

	_, idToValueMap := getNameAndIDRefMap(allVars, appID)

	for _, v := range allVars {
		if v.ApplicationID != appID {
			continue
		}
		appVars[v.Key] = v
		initText := v.TextValue
		switch true {
		case IsSystemVar(v.Key):
			continue
		case isSecret(v.Value):
			continue
		default:
			v.TextValue = replaceWithMap(v.Value, idToValueMap)
		}

		if initText != v.TextValue {
			updatedVars = append(updatedVars, v)
		}
	}
	return appVars, updatedVars
}
