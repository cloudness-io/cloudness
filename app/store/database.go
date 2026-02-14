package store

import (
	"context"
	"io"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

type (

	// InstanceStore defines the instance settings data storage.
	InstanceStore interface {
		// Create saves the instance settings.
		Create(ctx context.Context, instance *types.Instance) (*types.Instance, error)

		// Get gets the only instance settings
		Get(ctx context.Context) (*types.Instance, error)

		// Update updates the instance settings
		Update(ctx context.Context, instance *types.Instance) (*types.Instance, error)
	}

	// AuthSettingsStore defines the auth settings data storage
	AuthSettingsStore interface {
		// Create saves the auth setting.
		Create(ctx context.Context, auth *types.AuthSetting) (*types.AuthSetting, error)

		//Update updates the auth settings
		Update(ctx context.Context, auth *types.AuthSetting) (*types.AuthSetting, error)

		// FindByProvider gets the auth settings by provider.
		FindByProvider(ctx context.Context, provider enum.AuthProvider) (*types.AuthSetting, error)

		// List lists the auth settings.
		List(ctx context.Context) ([]*types.AuthSetting, error)
	}

	// ServerStore defines the server data storage
	ServerStore interface {
		//Find the server by id.
		Find(ctx context.Context, id int64) (*types.Server, error)

		//FindByUID the server by id.
		FindByUID(ctx context.Context, serverUID int64) (*types.Server, error)

		//Update  updates the server details
		Update(ctx context.Context, server *types.Server) (*types.Server, error)

		//Create save the server details
		Create(ctx context.Context, server *types.Server) (*types.Server, error)

		//List lists the servers
		List(ctx context.Context) ([]*types.Server, error)
	}

	// PrincipalStore defines the principal data storage.
	PrincipalStore interface {
		/*
		 * PRINCIPAL RELATED OPERATIONS.
		 */
		// Find finds the principal by id.
		Find(ctx context.Context, id int64) (*types.Principal, error)

		/*
		 * USER RELATED OPERATIONS.
		 */

		// FindUser finds the user by id.
		FindUser(ctx context.Context, id int64) (*types.User, error)

		// FindUserByUID finds the user by uid.
		FindUserByUID(ctx context.Context, uid string) (*types.User, error)

		// FindUserByEmail finds the user by email.
		FindUserByEmail(ctx context.Context, email string) (*types.User, error)

		// CreateUser saves the user details.
		CreateUser(ctx context.Context, user *types.User) (*types.User, error)

		// UpdateUser updates an existing user.
		UpdateUser(ctx context.Context, user *types.User) error

		// Count counts the user.
		CountUsers(ctx context.Context) (int64, error)
	}

	// TokenStore defines the token data storage.
	TokenStore interface {
		// Find finds the token by id
		Find(ctx context.Context, id int64) (*types.Token, error)

		// FindByIdentifier finds the token by principalId and token identifier.
		FindByIdentifier(ctx context.Context, principalID int64, identifier string) (*types.Token, error)

		// Create saves the token details.
		Create(ctx context.Context, token *types.Token) error

		// Delete deletes the token with the given id.
		Delete(ctx context.Context, id int64) error

		// DeleteExpiredBefore deletes all tokens that expired before the provided time.
		// If tokenTypes are provided, then only tokens of that type are deleted.
		DeleteExpiredBefore(ctx context.Context, before time.Time, tknTypes []enum.TokenType) (int64, error)

		// List returns a list of tokens of a specific type for a specific principal.
		List(ctx context.Context, principalID int64, tokenType enum.TokenType) ([]*types.Token, error)

		// Count returns a count of tokens of a specifc type for a specific principal.
		Count(ctx context.Context, principalID int64, tokenType enum.TokenType) (int64, error)
	}

	// TenantStore defines the tenent data storage
	TenantStore interface {
		//Find the tenant by id.
		Find(ctx context.Context, id int64) (*types.Tenant, error)

		//FindByUID the tenant by id.
		FindByUID(ctx context.Context, tenantUID int64) (*types.Tenant, error)

		//Update  updates the tenant details
		Update(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error)

		//Create save the tenant details
		Create(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error)

		//List lists the tenant by given filter
		List(ctx context.Context, filter *types.TenantFilter) ([]*types.Tenant, error)

		//SoftDelete deletes the tenant softly by setting the deleted timestamp
		SoftDelete(ctx context.Context, tenant *types.Tenant, deletedAt int64) error

		// Purge the soft deleted tenant permanently.
		Purge(ctx context.Context, id int64, deletedAt *int64) error
	}

	// TenantMembershipStore defines the tenant membership data storage
	TenantMembershipStore interface {
		//Get the membership by tenant id and principal id
		Find(ctx context.Context, tenantID, principalID int64) (*types.TenantMembership, error)

		//List the membership by principal id
		List(ctx context.Context, principalID int64) ([]*types.TenantMembershipUser, error)

		//List membership by tenant id
		ListByTenant(ctx context.Context, tenantID int64) ([]*types.TenantMembershipUser, error)

		//Create save the tenant membership
		Create(ctx context.Context, tenantMembership *types.TenantMembership) error

		//Update updates the tenant membership
		Update(ctx context.Context, tenantID, principalID int64, role enum.TenantRole) error

		//Delete deletes the membership
		Delete(ctx context.Context, tenantID, principalID int64) error

		//DeleteAll deletes all membership of a tenant
		DeleteAll(ctx context.Context, tenantID int64) error
	}

	// PrivateKeyStore defines the private key data storage
	PrivateKeyStore interface {
		// Find the private key by id.
		Find(ctx context.Context, tenantID, id int64) (*types.PrivateKey, error)

		//Create save the private key
		Create(ctx context.Context, privatekey *types.PrivateKey) (*types.PrivateKey, error)

		//Delete delets the private key
		Delete(ctx context.Context, tenantID, id int64) error
	}

	// ProjectStore defines the project data storage
	ProjectStore interface {
		// Find the project by id.
		Find(ctx context.Context, id int64) (*types.Project, error)

		//FindByUID the tenant by id.
		FindByUID(ctx context.Context, tenantID int64, projectUID int64) (*types.Project, error)

		//Count returns the number of projects for the given tenant
		Count(ctx context.Context, opts *types.ProjectFilter) (int64, error)

		//List lists the project by given filter
		List(ctx context.Context, opts *types.ProjectFilter) ([]*types.Project, error)

		// Create save the project details.
		Create(ctx context.Context, project *types.Project) (*types.Project, error)

		// Update updates the project details.
		Update(ctx context.Context, project *types.Project) (*types.Project, error)

		//SoftDelete deletes the project softly by setting the deleted timestamp
		SoftDelete(ctx context.Context, project *types.Project, deletedAt int64) error

		// Purge the soft deleted project permanently.
		Purge(ctx context.Context, id int64, deletedAt *int64) error
	}

	// ProjectMembershipStore defines the project membership data storage
	ProjectMembershipStore interface {
		//Get the membership by tenant id, project id and principal id
		Find(ctx context.Context, tenantID, projectID, principalID int64) (*types.ProjectMembership, error)

		// List the membership for a project
		List(ctx context.Context, tenantID, projectID int64) ([]*types.ProjectMembershipUser, error)

		// Create save the project membership
		Create(ctx context.Context, projectMembership *types.ProjectMembership) error

		//Update updates the project membership
		Update(ctx context.Context, tenantID, projectID, principalID int64, role enum.ProjectRole) error

		//Delete deletes the membership
		Delete(ctx context.Context, tenantID, projectID, principalID int64) error
	}

	// EnvironmentStore defines the project data storage
	EnvironmentStore interface {
		// Find the environment by id.
		Find(ctx context.Context, id int64) (*types.Environment, error)

		//FindByUID finds the environment by projectID and environmentUID.
		FindByUID(ctx context.Context, projectID int64, environmentUID int64) (*types.Environment, error)

		// List returns a list of environment for the given filter
		List(ctx context.Context, filter *types.EnvironmentFilter) ([]*types.Environment, error)

		// Create save the environment details.
		Create(ctx context.Context, environment *types.Environment) (*types.Environment, error)

		// Update updates the environment details.
		Update(ctx context.Context, environment *types.Environment) (*types.Environment, error)

		//SoftDelete deletes the environment softly by setting the deleted timestamp
		SoftDelete(ctx context.Context, environmment *types.Environment, deletedAt int64) error

		// Purge the soft deleted environment permanently.
		Purge(ctx context.Context, id int64, deletedAt *int64) error
	}

	// GithubAppStore defines the github app data storage.
	GithubAppStore interface {
		// Find the github app by id.
		Find(ctx context.Context, tenantID, projectID, githubAppID int64) (*types.GithubApp, error)

		//FindByUID finds the github app by UID.
		FindByUID(ctx context.Context, tenantID int64, projectID int64, githubAppUID int64) (*types.GithubApp, error)

		//List lists the github apps for tenant and project
		List(ctx context.Context, tenantID, projectID int64) ([]*types.GithubApp, error)

		// Update updates the github app.
		Update(ctx context.Context, githubapp *types.GithubApp) (*types.GithubApp, error)

		// Create save the github app.
		Create(ctx context.Context, githubapp *types.GithubApp) (*types.GithubApp, error)

		// Delete deletes the github app.
		Delete(ctx context.Context, githubapp *types.GithubApp) error
	}

	// ApplicationStore defines the application data storage
	ApplicationStore interface {
		// create save the application
		Create(ctx context.Context, application *types.Application) (*types.Application, error)

		// UpdateSpec updates application spec
		UpdateSpec(ctx context.Context, application *types.Application) (*types.Application, error)

		// UpdateDeploymentStatus updates application deployment status
		UpdateDeploymentStatus(ctx context.Context, application *types.Application) (*types.Application, error)

		// UpdateStatus updates application status
		UpdateStatus(ctx context.Context, appUID int64, status enum.ApplicationStatus) (*types.Application, error)

		// UpdateDeploymentTriggerTime updates the application to latest trigger time
		UpdateDeploymentTriggerTime(ctx context.Context, application *types.Application) (*types.Application, error)

		// UpdateNeedsDeployment updates the application needs deployment
		UpdateNeedsDeployment(ctx context.Context, application *types.Application) (*types.Application, error)

		// Find the application by id
		Find(ctx context.Context, id int64) (*types.Application, error)

		//FindByUID finds the application by tenant id, project id ,environment id and application u_id
		FindByUID(ctx context.Context, tenantID, projectID, environmentID, applicationUID int64) (*types.Application, error)

		//List lists the applications by filter
		List(ctx context.Context, filter *types.ApplicationFilter) ([]*types.Application, error)

		//Count counts the application by filter
		Count(ctx context.Context, filter *types.ApplicationFilter) (int64, error)

		//SoftDelete deletes the project softly by setting the deleted timestamp
		SoftDelete(ctx context.Context, application *types.Application, deletedAt int64) error

		// Purge the soft deleted application permanently.
		Purge(ctx context.Context, id int64, deletedAt *int64) error
	}

	VolumeStore interface {
		// Create creates a new volume
		Create(ctx context.Context, volume *types.Volume) (*types.Volume, error)

		// Update updates the volume
		Update(ctx context.Context, volume *types.Volume) (*types.Volume, error)

		// Find the volume by id
		Find(ctx context.Context, id int64) (*types.Volume, error)

		//FindByUID finds the volume by tenant id, project id ,environment id and volume u_id
		FindByUID(ctx context.Context, tenantID, projectID, environmentID, volumeUID int64) (*types.Volume, error)

		//List lists the applications by tenant id, project id and environment id
		List(ctx context.Context, filter *types.VolumeFilter) ([]*types.Volume, error)

		//SoftDelete deletes the volume softly by setting the deleted timestamp
		SoftDelete(ctx context.Context, volume *types.Volume, deletedAt int64) error

		// Purge the soft deleted volume permanently.
		Purge(ctx context.Context, id int64, deletedAt *int64) error
	}

	VariableStore interface {
		// Find the variable by application id and variable uid
		Find(ctx context.Context, applicationID, varUID int64) (*types.Variable, error)

		// Upsert updates or inserts the variable
		Upsert(ctx context.Context, variables *types.Variable) error

		// UpsertMany updates or inserts the variables
		UpsertMany(ctx context.Context, variables []*types.Variable) error

		// List lists the variables by environment id and application id
		List(ctx context.Context, environmentID, applicationID int64) ([]*types.Variable, error)

		// ListInEnvironment lists the variables by environment id
		ListInEnvironment(ctx context.Context, envID int64) ([]*types.Variable, error)

		// Delete deletes the variables by application id and varUID
		Delete(ctx context.Context, applicationID, varUID int64) error

		// DeleteByKey deletes the variables by application id and key
		DeleteByKey(ctx context.Context, applicationID int64, key string) error

		// DeleteByKeys deletes multiple variables by application id and keys
		DeleteByKeys(ctx context.Context, applicationID int64, keys []string) error
	}

	// DeploymentStore defines the deployment data storage
	DeploymentStore interface {
		// Find the deployment by id
		Find(ctx context.Context, id int64) (*types.Deployment, error)

		// FindByUID finds the deployment by deployment u_id and application id
		FindByUID(ctx context.Context, applicationID int64, deploymentUID int64) (*types.Deployment, error)

		// List lists the deployments by application id
		List(ctx context.Context, applicationID int64) ([]*types.Deployment, error)

		// Create save the deployment
		Create(ctx context.Context, deployment *types.Deployment) (*types.Deployment, error)

		// Update updates a deployment
		Update(ctx context.Context, deployment *types.Deployment) error

		// ListIncomplete lists the incomplete deployments
		ListIncomplete(ctx context.Context) ([]*types.Deployment, error)

		// ListIncompleteByApplicationID lists the incomplete deployments for a application
		ListIncompleteByApplicationID(ctx context.Context, applicationID int64) ([]*types.Deployment, error)
	}

	// LogStore defines the log data storage
	LogStore interface {
		// Find the log by deployment id
		Find(ctx context.Context, deploymentID int64) (io.ReadCloser, error)

		// Create writes copies of log stream from reader to store
		Create(ctx context.Context, deploymentID int64, r io.Reader) error
	}

	// TemplateStore defines the template data storage
	TemplateStore interface {
		//Find the template by id
		Find(ctx context.Context, id int64) (*types.Template, error)

		// UpsertMany updates or inserts the templates
		UpsertMany(ctx context.Context, templates []*types.Template) error

		// List lists the templates
		List(ctx context.Context) ([]*types.Template, error)

		// ListTags lists all template tags
		ListTags(ctx context.Context) ([]*types.Tag, error)

		// ListByTag lists templates associated with a tag slug
		ListByTag(ctx context.Context, tag string) ([]*types.Template, error)

		// List by slugs
		ListBySlugs(ctx context.Context, slugs []string) ([]*types.Template, error)

		// List not in slugs
		ListNotInSlugs(ctx context.Context, slugs []string) ([]*types.Template, error)
	}

	// FavoriteStore defines the favorite data storage
	FavoriteStore interface {
		// List lists the favorite by user id and tenant id
		List(ctx context.Context, userID, tenantID int64) ([]*types.FavoriteDTO, error)

		// Find find the favorite by userid and application id
		Find(ctx context.Context, userID, applicationID int64) (*types.Favorite, error)

		// Add adds a favorite
		Add(ctx context.Context, userID, applicationID int64) error

		// Delete delete a favorite
		Delete(ctx context.Context, userID, applicationID int64) error
	}

	// MetricsStore defines the metrics data storage
	MetricsStore interface {
		// UpsertMany upserts the metrics
		UpsertMany(ctx context.Context, metrics []*types.AppMetrics) error

		// ListByApplicationUID lists the metrics by application uid
		ListByApplicationUID(ctx context.Context, applicationUID int64, from time.Time, to time.Time, bucketSeconds int64) ([]*types.AppMetricsAggregate, error)
	}
)
