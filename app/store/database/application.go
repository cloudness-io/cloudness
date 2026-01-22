package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.ApplicationStore = (*ApplicationStore)(nil)

// NewApplicationStore returns new ApplicationStore
func NewApplicationStore(db *sqlx.DB) *ApplicationStore {
	return &ApplicationStore{
		db: db,
	}
}

// ApplicationStore implements a ApplicationStore backed by a relational database.
type ApplicationStore struct {
	db *sqlx.DB
}

// application is a DB representation of a application.
type application struct {
	types.Application
}

// applicationColumns defines the columns that are the same across all applications.
const applicationColumns = `
	application_id
	,application_uid
	,application_tenant_id
	,application_project_id
	,application_environment_id
	,application_environment_uid
	,application_server_id
	,application_name
	,application_description
	,application_type
	,application_status
	,application_spec
	,application_githubapp_id
	,application_domain
	,application_custom_domain
	,application_private_domain
	,application_deployment_id
	,application_deployment_status
	,application_deployment_triggered_at
	,application_created
	,application_updated
	,application_deleted`

const applicationInsert = `
INSERT INTO applications  (
	application_uid
	,application_tenant_id
	,application_project_id
	,application_environment_id
	,application_environment_uid
	,application_server_id
	,application_name
	,application_description
	,application_type
	,application_status
	,application_spec
	,application_githubapp_id
	,application_domain
	,application_custom_domain
	,application_private_domain
	,application_deployment_id
	,application_deployment_status
	,application_deployment_triggered_at
	,application_created
	,application_updated	
	,application_deleted
) values (
	:application_uid
	,:application_tenant_id
	,:application_project_id
	,:application_environment_id
	,:application_environment_uid
	,:application_server_id
	,:application_name
	,:application_description
	,:application_type
	,:application_status
	,:application_spec
	,:application_githubapp_id
	,:application_domain
	,:application_custom_domain
	,:application_private_domain
	,:application_deployment_id
	,:application_deployment_status
	,:application_deployment_triggered_at
	,:application_created
	,:application_updated
	,:application_deleted
	) RETURNING application_id
	`

func (s *ApplicationStore) Create(ctx context.Context, application *types.Application) (*types.Application, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(applicationInsert, application)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind application object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&application.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert application query failed")
	}
	return application, nil
}

func (s *ApplicationStore) UpdateSpec(ctx context.Context, application *types.Application) (*types.Application, error) {
	const applicationUpdate = `UPDATE applications 
		SET
			application_name = :application_name
 			,application_description = :application_description
			,application_type = :application_type
 			,application_spec = :application_spec
			,application_domain = :application_domain
			,application_custom_domain = :application_custom_domain
			,application_private_domain = :application_private_domain
			,application_updated = :application_updated
			,application_deleted = :application_deleted
		WHERE application_id = :application_id`

	return s.update(ctx, application, applicationUpdate)
}

func (s *ApplicationStore) UpdateDeploymentStatus(ctx context.Context, application *types.Application) (*types.Application, error) {
	const applicationUpdateDeploymentStatus = `UPDATE applications 
		SET
			application_deployment_id = :application_deployment_id
			,application_deployment_status = :application_deployment_status
		WHERE application_id = :application_id`

	return s.update(ctx, application, applicationUpdateDeploymentStatus)
}

func (s *ApplicationStore) UpdateStatus(ctx context.Context, application *types.Application) (*types.Application, error) {
	const applicationUpdateStatus = `UPDATE applications 
		SET
			application_status = :application_status
		WHERE application_id = :application_id`

	return s.update(ctx, application, applicationUpdateStatus)
}

func (s *ApplicationStore) UpdateDeploymentTriggerTime(ctx context.Context, application *types.Application) (*types.Application, error) {
	const applicationUpdateTriggerTime = `UPDATE applications 
		SET
			application_deployment_triggered_at = :application_deployment_triggered_at
		WHERE application_id = :application_id`

	application.DeploymentTriggeredAt = time.Now().UTC().UnixMilli()
	return s.update(ctx, application, applicationUpdateTriggerTime)
}

func (s *ApplicationStore) UpdateNeedsDeployment(ctx context.Context, application *types.Application) (*types.Application, error) {
	const applicationUpdateTriggerTime = `UPDATE applications 
		SET
			application_updated = :application_updated
		WHERE application_id = :application_id`

	application.Updated = time.Now().UTC().UnixMilli()
	return s.update(ctx, application, applicationUpdateTriggerTime)
}

func (s *ApplicationStore) update(ctx context.Context, application *types.Application, query string) (*types.Application, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(query, application)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind application object")
	}

	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update application query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return application, nil
}

// Find the app by id.
func (s *ApplicationStore) Find(ctx context.Context, id int64) (*types.Application, error) {
	return s.FindDeleted(ctx, id, nil)
}

func (s *ApplicationStore) FindDeleted(ctx context.Context, id int64, deletedAt *int64) (*types.Application, error) {
	stms := database.Builder.
		Select(applicationColumns).
		From("applications").
		Where("application_id = ?", id)

	if deletedAt != nil {
		stms = stms.Where("application_deleted = ?", deletedAt)
	} else {
		stms = stms.Where("application_deleted IS NULL")
	}
	sql, args, err := stms.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(application)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDbApplication(dst)
}

func (s *ApplicationStore) FindByUID(ctx context.Context, tenantID, projectID, environmentID, uid int64) (*types.Application, error) {
	return s.FindDeletedByUID(ctx, tenantID, projectID, environmentID, uid, nil)
}

func (s *ApplicationStore) FindDeletedByUID(ctx context.Context, tenantID, projectID, environmentID, uid int64, deletedAt *int64) (*types.Application, error) {
	stmt := database.Builder.
		Select(applicationColumns).
		From("applications").
		Where("application_tenant_id = ?", tenantID).
		Where("application_project_id = ?", projectID).
		Where("application_environment_id = ?", environmentID).
		Where("application_uid = ?", uid)

	if deletedAt != nil {
		stmt = stmt.Where("application_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("application_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(application)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by uid query failed")
	}
	return s.mapDbApplication(dst)
}

// Find is the app name available.
func (s *ApplicationStore) CountByName(ctx context.Context, name string) (int64, error) {
	stmt := database.Builder.
		Select("count(1)").
		From("applications")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed executing count query")
	}
	return count, nil
}

// List returns a list of applications for the given filter
func (s *ApplicationStore) List(ctx context.Context, filter *types.ApplicationFilter) ([]*types.Application, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*application{}

	stmt := database.Builder.
		Select(applicationColumns).
		From("applications")

	stmt = s.applyQueryFilter(stmt, filter)
	stmt = s.applySortFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by parent id query failed")

	}
	return s.mapDBApplications(dst)
}

func (s *ApplicationStore) Count(ctx context.Context, filter *types.ApplicationFilter) (int64, error) {
	stmt := database.Builder.
		Select("count(1)").
		Where("application_deleted IS NULL").
		From("applications")

	stmt = s.applyQueryFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed executing count query")
	}
	return count, nil
}

func (s *ApplicationStore) SoftDelete(ctx context.Context, application *types.Application, deletedAt int64) error {
	application.Deleted = &deletedAt
	_, err := s.UpdateSpec(ctx, application)
	return err
}

// Purge deletes the application permanently.
func (s *ApplicationStore) Purge(ctx context.Context, id int64, deletedAt *int64) error {
	stmt := database.Builder.
		Delete("applications").
		Where("application_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("application_deleted = ?", deletedAt)
	} else {
		stmt = stmt.Where("application_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to delete application")
	}
	return nil
}

func (s *ApplicationStore) applyQueryFilter(stmt squirrel.SelectBuilder, filter *types.ApplicationFilter) squirrel.SelectBuilder {
	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("application_name", filter.Query))
	}
	//nolint:gocritic
	if filter.DeletedAt != nil {
		stmt = stmt.Where("application_deleted = ?", filter.DeletedAt)
	} else if filter.DeletedBeforeOrAt != nil {
		stmt = stmt.Where("application_deleted <= ?", filter.DeletedBeforeOrAt)
	} else {
		stmt = stmt.Where("application_deleted IS NULL")
	}

	if filter.TenantID != nil {
		stmt = stmt.Where("application_tenant_id = ?", filter.TenantID)
	}

	if filter.ProjectID != nil {
		stmt = stmt.Where("application_project_id = ?", filter.ProjectID)
	}

	if filter.EnvironmentID != nil {
		stmt = stmt.Where("application_environment_id = ?", filter.EnvironmentID)
	}

	return stmt
}

func (s *ApplicationStore) applySortFilter(stmt squirrel.SelectBuilder, filter *types.ApplicationFilter) squirrel.SelectBuilder {
	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	switch filter.Sort {
	case enum.ApplicationAttrNone:
		stmt = stmt.OrderBy("application_id " + enum.OrderAsc.String())
	case enum.ApplicationAttrCreated:
		stmt = stmt.OrderBy("application_created " + filter.Order.String())
	case enum.ApplicationAttrUpdated:
		stmt = stmt.OrderBy("application_updated" + filter.Order.String())
	case enum.ApplicationAttrDeleted:
		stmt = stmt.OrderBy("application_deleted" + filter.Order.String())
	}

	return stmt
}

func (s *ApplicationStore) mapDBApplications(dbApp []*application) ([]*types.Application, error) {
	res := make([]*types.Application, len(dbApp))
	for i := range dbApp {
		c, err := s.mapDbApplication(dbApp[i])
		if err != nil {
			return nil, err
		}
		res[i] = c
	}
	return res, nil
}

func (s *ApplicationStore) mapDbApplication(dbApplication *application) (*types.Application, error) {
	app := &dbApplication.Application

	compSpec := new(types.ApplicationSpec)
	if err := json.Unmarshal([]byte(app.SpecJSON), &compSpec); err != nil {
		return nil, err
	}
	app.Spec = compSpec
	return app, nil
}
