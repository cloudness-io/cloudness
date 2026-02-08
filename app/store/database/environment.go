package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/helpers"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.EnvironmentStore = (*EnvironmentStore)(nil)

func NewEnvironmentSore(db *sqlx.DB) *EnvironmentStore {
	return &EnvironmentStore{
		db: db,
	}
}

type EnvironmentStore struct {
	db *sqlx.DB
}

type environment struct {
	types.Environment
}

const environmentColumns = `
   environment_id,
   environment_uid,
	environment_tenant_id,
   environment_project_id,
	environment_sequence,
   environment_name,
   environment_created_by,
   environment_created,
   environment_updated,
	environment_deleted`

const environmentInsert = `
INSERT INTO environments (
   environment_uid
	,environment_tenant_id
   ,environment_project_id
	,environment_sequence
   ,environment_name
   ,environment_created_by
   ,environment_created
   ,environment_updated
) values (
   :environment_uid
	,:environment_tenant_id
   ,:environment_project_id
	,:environment_sequence
   ,:environment_name
   ,:environment_created_by
   ,:environment_created
   ,:environment_updated
) RETURNING environment_id`

const environmentUpdate = `
   UPDATE environments
   SET
	environment_sequence = :environment_sequence,
   environment_name = :environment_name,
   environment_updated = :environment_updated,
	environment_deleted = :environment_deleted
   WHERE environment_project_id = :environment_project_id AND environment_id = :environment_id
`

// Find the environment by id.
func (s *EnvironmentStore) Find(ctx context.Context, id int64) (*types.Environment, error) {
	return s.FindDeleted(ctx, id, nil)
}

func (s *EnvironmentStore) FindDeleted(ctx context.Context, id int64, deletedAt *int64) (*types.Environment, error) {
	stmt := database.Builder.
		Select(environmentColumns).
		From("environments").
		Where("environment_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("environment_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("environment_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(environment)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDBEnvironment(dst), nil
}

// FindByUID finds the environment by projectID and environmentUID.
func (s *EnvironmentStore) FindByUID(ctx context.Context, projectID int64, environmentUID int64) (*types.Environment, error) {
	return s.FindDeletedByUID(ctx, projectID, environmentUID, nil)
}

func (s *EnvironmentStore) FindDeletedByUID(ctx context.Context, projectID int64, environmentUID int64, deletedAt **int64) (*types.Environment, error) {
	stmt := database.Builder.
		Select(environmentColumns).
		From("environments").
		Where("environment_project_id = ?", projectID).
		Where("environment_uid = ?", environmentUID)

	if deletedAt != nil {
		stmt = stmt.Where("environment_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("environment_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(environment)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select environment by project id and environment uid query failed")
	}
	return s.mapDBEnvironment(dst), nil
}

// List returns a list of environment for the given filter
func (s *EnvironmentStore) List(ctx context.Context, filter *types.EnvironmentFilter) ([]*types.Environment, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*environment{}

	stmt := database.Builder.
		Select(environmentColumns).
		From("environments")

	stmt = s.applyQueryFilter(stmt, filter)
	stmt = s.applySortFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by parent id query failed")

	}
	return s.mapDBEnvironmentss(dst)
}

// Create save the environment details.
func (s *EnvironmentStore) Create(ctx context.Context, environment *types.Environment) (*types.Environment, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(environmentInsert, environment)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind environment object")
	}

	if err = db.QueryRowxContext(ctx, query, arg...).Scan(&environment.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert environment query failed")
	}
	return environment, nil

}

// Update updates the environment details.
func (s *EnvironmentStore) Update(ctx context.Context, environment *types.Environment) (*types.Environment, error) {
	environment.Updated = time.Now().UTC().UnixMilli()
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(environmentUpdate, environment)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind environment object")
	}

	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update environment query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return environment, nil
}

func (s *EnvironmentStore) SoftDelete(ctx context.Context, environment *types.Environment, deletedAt int64) error {
	environment.Deleted = &deletedAt
	environment.Seq = helpers.RandomNum(9000, 9999)
	_, err := s.Update(ctx, environment)
	return err
}

func (s *EnvironmentStore) Purge(ctx context.Context, id int64, deletedAt *int64) error {
	stmt := database.Builder.
		Delete("environments").
		Where("environment_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("environment_deleted = ?", deletedAt)
	} else {
		stmt = stmt.Where("environment_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrapf(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to delete environment")
	}
	return nil
}

func (s *EnvironmentStore) applyQueryFilter(stmt squirrel.SelectBuilder, filter *types.EnvironmentFilter) squirrel.SelectBuilder {
	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("environment_name", filter.Query))
	}

	if filter.TenantID != nil {
		stmt = stmt.Where("environment_tenant_id = ?", filter.TenantID)
	}
	if filter.ProjectID != nil {
		stmt = stmt.Where("environment_project_id = ?", filter.ProjectID)
	}

	//nolint:gocritic
	if filter.DeletedAt != nil {
		stmt = stmt.Where("environment_deleted = ?", filter.DeletedAt)
	} else if filter.DeletedBeforeOrAt != nil {
		stmt = stmt.Where("environment_deleted <= ?", filter.DeletedBeforeOrAt)
	} else {
		stmt = stmt.Where("environment_deleted IS NULL")
	}

	return stmt
}

func (s *EnvironmentStore) applySortFilter(stmt squirrel.SelectBuilder, filter *types.EnvironmentFilter) squirrel.SelectBuilder {
	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	switch filter.Sort {
	case enum.EnvironmentAttrNone:
		stmt = stmt.OrderBy("environment_id " + enum.OrderAsc.String())
	case enum.EnvironmentAttrCreated:
		stmt = stmt.OrderBy("environment_created " + filter.Order.String())
	case enum.EnvironmentAttrUpdated:
		stmt = stmt.OrderBy("environment_updated " + filter.Order.String())
	case enum.EnvironmentAttrDeleted:
		stmt = stmt.OrderBy("environment_deleted " + filter.Order.String())
	case enum.EnvironmentAttrSequence:
		stmt = stmt.OrderBy("environment_sequence " + filter.Order.String())
	}

	return stmt
}

func (s *EnvironmentStore) mapDBEnvironment(src *environment) *types.Environment {
	return &src.Environment
}

func (s *EnvironmentStore) mapDBEnvironmentss(src []*environment) ([]*types.Environment, error) {
	res := make([]*types.Environment, len(src))
	for i := range src {
		res[i] = s.mapDBEnvironment(src[i])
	}
	return res, nil
}
