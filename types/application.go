package types

import (
	"encoding/json"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Application struct {
	ID                    int64                            `db:"application_id"                           json:"-"`
	UID                   int64                            `db:"application_uid"                          json:"uid"`
	TenantID              int64                            `db:"application_tenant_id"                    json:"-"`
	ProjectID             int64                            `db:"application_project_id"                   json:"-"`
	EnvironmentID         int64                            `db:"application_environment_id"               json:"-"`
	EnvironmentUID        int64                            `db:"application_environment_uid"              json:"environment_uid"`
	ServerID              int64                            `db:"application_server_id"                    json:"-"`
	Name                  string                           `db:"application_name"                         json:"name"`
	Slug                  string                           `db:"application_slug"                         json:"slug"`
	ParentSlug            string                           `db:"application_parent_slug"                  json:"parent_slug"`
	Description           string                           `db:"application_description"                  json:"description"`
	Type                  enum.ApplicationType             `db:"application_type"                         json:"type"`
	Status                enum.ApplicationStatus           `db:"application_status"                       json:"status"`
	SpecJSON              string                           `db:"application_spec"                         json:"-"`
	Spec                  *ApplicationSpec                 `db:"-"                                        json:"spec"`
	GithubAppID           *int64                           `db:"application_githubapp_id"                 json:"-"`
	Domain                string                           `db:"application_domain"                       json:"domain"`
	CustomDomain          string                           `db:"application_custom_domain"                json:"custom_domain"`
	PrivateDomain         string                           `db:"application_private_domain"               json:"private_domain"`
	DeploymentID          *int64                           `db:"application_deployment_id"                json:"-"`
	DeploymentStatus      enum.ApplicationDeploymentStatus `db:"application_deployment_status"            json:"deployment_status"`
	DeploymentTriggeredAt int64                            `db:"application_deployment_triggered_at"      json:"deployment_triggered_at"`

	Created int64  `db:"application_created"                json:"created"`
	Updated int64  `db:"application_updated"                json:"updated"`
	Deleted *int64 `db:"application_deleted"                json:"deleted"`
}

type ApplicationFilter struct {
	ListQueryFilter
	TenantID          *int64               `json:"tenant_id,omitempty"`
	ProjectID         *int64               `json:"project_id,omitempty"`
	EnvironmentID     *int64               `json:"environment_id,omitempty"`
	Sort              enum.ApplicationAttr `json:"sort"`
	Order             enum.Order           `json:"order"`
	DeletedAt         *int64               `json:"deleted_at,omitempty"`
	DeletedBeforeOrAt *int64               `json:"deleted_before_or_at,omitempty"`
}

// helpers
func (a *Application) GetIdentifierStr() string {
	return a.Slug
}

// GetGithubAppID returns the github app id if it exists
func (a *Application) GetGithubAppID() int64 {
	if a == nil || a.GithubAppID == nil {
		return 0
	}
	return *a.GithubAppID
}

func (a *Application) UpdateSpecJSON() error {
	specJSON, err := json.Marshal(a.Spec)
	if err != nil {
		return errors.BadRequest("Error converting spec to JSON")
	}
	a.SpecJSON = string(specJSON)
	return nil
}

func (a *Application) IsStateful() bool {
	return a.Type == enum.ApplicationTypeStateful
}

func (a *Application) IsStateless() bool {
	return a.Type == enum.ApplicationTypeStateless
}
