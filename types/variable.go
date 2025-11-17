package types

import (
	"github.com/cloudness-io/cloudness/types/enum"
)

type Variable struct {
	UID             int64             `db:"variable_uid"              json:"uid"`
	EnvironmentID   int64             `db:"variable_environment_id"   json:"-"`
	ApplicationID   int64             `db:"variable_application_id"   json:"-"`
	ApplicationName string            `db:"application_name"          json:"application_name"`
	Key             string            `db:"variable_key"              json:"key"`
	Value           string            `db:"variable_value"            json:"value"`
	TextValue       string            `db:"variable_text_value"       json:"text_value"`
	Type            enum.VariableType `db:"variable_type"             json:"type"`
	Created         int64             `db:"variable_created"          json:"created"`
	Updated         int64             `db:"variable_updated"          json:"updated"`
}

// DTO to render frontend
type VariableDTO struct {
	SystemVariable []*SystemVariable
	UserVariable   []*UserVariable
}

type SystemVariable struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	TextValue string `json:"text_value"`
}

type UserVariable struct {
	UID            int64  `json:"uid"`
	Key            string `json:"key"`
	Type           string `json:"type"`
	TextValue      string `json:"text_value"`
	ReferenceValue string `json:"reference_value"`
	IsGenerated    bool   `json:"is_generated"`
}
