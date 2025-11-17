package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.GithubAppStore = (*GithubAppStore)(nil)

type GithubAppStore struct {
	db *sqlx.DB
}

func NewGithubAppStore(db *sqlx.DB) *GithubAppStore {
	return &GithubAppStore{
		db: db,
	}
}

// githubapp is a DB representation of a SourceProvider.
type githubapp struct {
	types.GithubApp
}

// githubappColumns defines the columns that are the same accross all integrations.
const githubappColumns = `
	github_app_id
	,github_app_uid
	,github_app_tenant_id
	,github_app_project_id
	,github_app_private_key_id
	,github_app_is_tenant_wide
	,github_app_name
	,github_app_organization
	,github_app_api_url
	,github_app_html_url
	,github_app_custom_user
	,github_app_custom_port
	,github_app_app_id
	,github_app_installation_id
	,github_app_client_id
	,github_app_client_secret
	,github_app_webhook_secret
   ,github_app_created_by
   ,github_app_created
   ,github_app_updated
`

const githubappInsert = `
INSERT INTO github_apps (
	github_app_uid
	,github_app_tenant_id
	,github_app_project_id
	,github_app_private_key_id
	,github_app_is_tenant_wide
	,github_app_name
	,github_app_organization
	,github_app_api_url
	,github_app_html_url
	,github_app_custom_user
	,github_app_custom_port
	,github_app_app_id
	,github_app_installation_id
	,github_app_client_id
	,github_app_client_secret
	,github_app_webhook_secret
   ,github_app_created_by
   ,github_app_created
   ,github_app_updated
) values (
	:github_app_uid
	,:github_app_tenant_id
	,:github_app_project_id
	,:github_app_private_key_id
	,:github_app_is_tenant_wide
	,:github_app_name
	,:github_app_organization
	,:github_app_api_url
	,:github_app_html_url
	,:github_app_custom_user
	,:github_app_custom_port
	,:github_app_app_id
	,:github_app_installation_id
	,:github_app_client_id
	,:github_app_client_secret
	,:github_app_webhook_secret
	,:github_app_created_by
	,:github_app_created
	,:github_app_updated
	) RETURNING github_app_id`

const githubappUpdate = `UPDATE github_apps 
	SET 
		github_app_name = :github_app_name
		,github_app_private_key_id = :github_app_private_key_id
		,github_app_app_id = :github_app_app_id
		,github_app_installation_id = :github_app_installation_id
		,github_app_client_id = :github_app_client_id
		,github_app_client_secret = :github_app_client_secret
		,github_app_webhook_secret = :github_app_webhook_secret
		,github_app_updated = :github_app_updated
	WHERE github_app_id = :github_app_id AND github_app_tenant_id = :github_app_tenant_id`

const githubappSelectBase = `
SELECT` + githubappColumns + `
FROM github_apps`

func (s *GithubAppStore) Find(ctx context.Context, tenantID, projectID, githubAppID int64) (*types.GithubApp, error) {
	const sqlQuery = githubappSelectBase + `
	WHERE github_app_tenant_id =$1 AND github_app_project_id = $2 AND github_app_id = $3`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(githubapp)

	if err := db.GetContext(ctx, dst, sqlQuery, tenantID, projectID, githubAppID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBGithubApp(dst), nil
}

func (s *GithubAppStore) FindByUID(ctx context.Context, tenantID int64, projectID int64, githubAppUID int64) (*types.GithubApp, error) {
	const sqlQuery = githubappSelectBase + `
	WHERE github_app_tenant_id =$1 AND github_app_uid = $2`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(githubapp)
	if err := db.GetContext(ctx, dst, sqlQuery, tenantID, githubAppUID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "select github app by uid query failed")
	}

	if !dst.IsTenantWide && dst.ProjectID != projectID {
		return nil, dbStore.ErrResourceNotFound
	}

	return s.mapDBGithubApp(dst), nil
}

// List lists the github apps for tenant and project
func (s *GithubAppStore) List(ctx context.Context, tenantID, projectID int64) ([]*types.GithubApp, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*githubapp{}

	stmt := database.Builder.
		Select(githubappColumns).
		From("github_apps").
		Where("github_app_tenant_id = $1 AND (github_app_project_id = $2 OR github_app_is_tenant_wide = true)")

	sql, _, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, tenantID, projectID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List github app by tenant and project query failed")
	}

	return s.mapDBGithubApps(dst), nil
}

// Update updates the github app.
func (s *GithubAppStore) Update(ctx context.Context, githubapp *types.GithubApp) (*types.GithubApp, error) {
	githubapp.Updated = time.Now().UTC().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(githubappUpdate, githubapp)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind github app object")
	}

	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update github app query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, dbStore.ErrVersionConflict
	}

	return githubapp, nil
}

// Create save the github app.
func (s *GithubAppStore) Create(ctx context.Context, githubapp *types.GithubApp) (*types.GithubApp, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(githubappInsert, githubapp)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind application object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&githubapp.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert application query failed")
	}
	return githubapp, nil
}

func (s *GithubAppStore) Delete(ctx context.Context, githubapp *types.GithubApp) error {
	githubappDeleteStmt := `DELETE FROM github_apps WHERE github_app_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, githubappDeleteStmt, githubapp.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete github app query failed")
	}
	return nil
}

func (s *GithubAppStore) mapDBGithubApps(src []*githubapp) []*types.GithubApp {
	var dst []*types.GithubApp
	for _, v := range src {
		dst = append(dst, s.mapDBGithubApp(v))
	}
	return dst
}

func (s *GithubAppStore) mapDBGithubApp(src *githubapp) *types.GithubApp {
	return &src.GithubApp
}
