package types

type PrivateKey struct {
	ID          int64  `db:"private_key_id"           json:"-"`
	UID         int64  `db:"private_key_uid"          json:"uid"`
	TenantID    int64  `db:"private_key_tenant_id"    json:"-"`
	Name        string `db:"private_key_name"         json:"name"`
	Description string `db:"private_key_description"  json:"description"`
	Key         string `db:"private_key_pem"          json:"-"`
	IsGit       bool   `db:"private_key_is_git"       json:"is_git"`

	Created int64 `db:"private_key_created"                json:"created"`
	Updated int64 `db:"private_key_updated"                json:"updated"`
}
