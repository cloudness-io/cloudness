package types

import "github.com/cloudness-io/cloudness/types/enum"

type AuthSetting struct {
	ID           int64             `db:"auth_id"              json:"-"`
	Provider     enum.AuthProvider `db:"auth_provider"        json:"provider"`
	Enabled      bool              `db:"auth_enabled"         json:"enabled,string"`
	ClientID     string            `db:"auth_client_id"       json:"client_id"`
	ClientSecret string            `db:"auth_client_secret"   json:"client_secret"`
	BaseURL      string            `db:"auth_base_url"        json:"base_url"`
	Created      int64             `db:"auth_created"         json:"created"`
	Updated      int64             `db:"auth_updated"         json:"updated"`
}
