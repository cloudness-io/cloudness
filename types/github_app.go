package types

import "fmt"

type GithubApp struct {
	ID             int64  `db:"github_app_id"                json:"-"`
	UID            int64  `db:"github_app_uid"               json:"uid"`
	TenantID       int64  `db:"github_app_tenant_id"         json:"-"`
	ProjectID      int64  `db:"github_app_project_id"        json:"-"`
	PrivateKeyID   int64  `db:"github_app_private_key_id"    json:"private_key_id"`
	IsTenantWide   bool   `db:"github_app_is_tenant_wide"    json:"is_tenant_wide"`
	Name           string `db:"github_app_name"              json:"name"`
	Organization   string `db:"github_app_organization"      json:"organization"`
	ApiUrl         string `db:"github_app_api_url"           json:"api_url"`
	HtmlUrl        string `db:"github_app_html_url"          json:"html_url"`
	CustomUser     string `db:"github_app_custom_user"       json:"custom_user"`
	CustomPort     int64  `db:"github_app_custom_port"       json:"custom_port"`
	AppID          int64  `db:"github_app_app_id"            json:"app_id"`
	InstallationID int64  `db:"github_app_installation_id"   json:"installation_id"`
	ClientID       string `db:"github_app_client_id"         json:"client_id"`
	ClientSecret   string `db:"github_app_client_secret"     json:"client_secret"`
	WebhookSecret  string `db:"github_app_webhook_secret"    json:"webhook_secret"`
	CreateBy       int64  `db:"github_app_created_by"        json:"-"`

	Created int64 `db:"github_app_created"                json:"created"`
	Updated int64 `db:"github_app_updated"                json:"updated"`
}

// Helpers
func (g *GithubApp) GetHttpUrl(repo string) string {
	return fmt.Sprintf("%s/%s.git", g.HtmlUrl, repo)
}
