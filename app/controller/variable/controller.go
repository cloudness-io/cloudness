package variable

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Controller struct {
	variableStore store.VariableStore
}

func NewController(variableStore store.VariableStore) *Controller {
	return &Controller{
		variableStore: variableStore,
	}
}

type AddVariableInput struct {
	Key   string            `json:"new_key"`
	Value string            `json:"new_value"`
	Type  enum.VariableType `json:"type"`
}
