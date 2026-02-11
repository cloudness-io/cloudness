package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.ProjectStore = (*ProjectStore)(nil)

func NewProjectSore(db *sqlx.DB) *ProjectStore {
	return &ProjectStore{
		db: db,
	}
}

type ProjectStore struct {
	db *sqlx.DB
}

type project struct {
	types.Project
}

const projectColumns = `
	project_id,
	project_tenant_id,
	project_name,
	project_uid,
	project_slug,
	project_description,
	project_created_by,
	project_created,
	project_updated,
	project_deleted`

const projectInsert = `
INSERT INTO projects (
	project_uid
	,project_slug
	,project_tenant_id
	,project_name
	,project_description
	,project_created_by
	,project_created
	,project_updated
	,project_deleted
) values (
	:project_uid
	,:project_slug
	,:project_tenant_id
	,:project_name
	,:project_description
	,:project_created_by
	,:project_created
	,:project_updated
	,:project_deleted
	) RETURNING project_id
	`

const projectUpdate = `UPDATE projects
	SET 
		project_name = :project_name
		,project_description = :project_description
		,project_updated = :project_updated
		,project_deleted = :project_deleted
	WHERE project_id = :project_id and project_tenant_id = :project_tenant_id`

func (s *ProjectStore) Find(ctx context.Context, id int64) (*types.Project, error) {
	return s.FindDeleted(ctx, id, nil)
}

func (s *ProjectStore) FindDeleted(ctx context.Context, id int64, deletedAt *int64) (*types.Project, error) {
	stmt := database.Builder.
		Select(projectColumns).
		From("projects").
		Where("project_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("project_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("project_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(project)
	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDBProject(dst), nil
}

func (s *ProjectStore) FindByUID(ctx context.Context, tenantID int64, projectUID int64) (*types.Project, error) {
	return s.FindDeletedByUID(ctx, tenantID, projectUID, nil)
}

func (s *ProjectStore) FindDeletedByUID(ctx context.Context, tenantID int64, projectUID int64, deletedAt **int64) (*types.Project, error) {
	stmt := database.Builder.
		Select(projectColumns).
		From("projects").
		Where("project_tenant_id = ?", tenantID).
		Where("project_uid = ?", projectUID)

	if deletedAt != nil {
		stmt = stmt.Where("project_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("project_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(project)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select project by uid query failed")
	}
	return s.mapDBProject(dst), nil
}

func (s *ProjectStore) Count(ctx context.Context, filter *types.ProjectFilter) (int64, error) {
	stmt := database.Builder.
		Select("count(1)").
		From("projects")

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

func (s *ProjectStore) List(ctx context.Context, filter *types.ProjectFilter) ([]*types.Project, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*project{}

	stmt := database.Builder.
		Select(projectColumns).
		From("projects")

	if filter.PrincipalID != nil {
		stmt = stmt.
			LeftJoin(`project_memberships 
				ON project_memberships.project_membership_project_id = projects.project_id
				AND project_memberships.project_membership_principal_id = ?`, filter.PrincipalID).
			LeftJoin(`tenant_memberships 
				ON tenant_memberships.tenant_membership_tenant_id = projects.project_tenant_id 
				AND tenant_memberships.tenant_membership_principal_id = ?`, filter.PrincipalID).
			Where(
				sq.Or{
					sq.Eq{`tenant_memberships.tenant_membership_role`: enum.TenantRoleAdmin},
					sq.NotEq{`project_memberships.project_membership_role`: nil},
				},
			)
	}

	stmt = s.applyQueryFilter(stmt, filter)
	stmt = s.applySortFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by parent id query failed")

	}
	return s.mapDBProjects(dst), nil
}

func (s *ProjectStore) Create(ctx context.Context, project *types.Project) (*types.Project, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(projectInsert, project)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind project object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&project.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert project query failed")
	}

	return project, nil
}

func (s *ProjectStore) Update(ctx context.Context, project *types.Project) (*types.Project, error) {
	project.Updated = time.Now().UTC().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(projectUpdate, project)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind project object")
	}
	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update project query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return project, nil
}

func (s *ProjectStore) SoftDelete(ctx context.Context, project *types.Project, deletedAt int64) error {
	project.Deleted = &deletedAt

	_, err := s.Update(ctx, project)
	return err
}

// Purge deletes the project permanently.
func (s *ProjectStore) Purge(ctx context.Context, id int64, deletedAt *int64) error {
	stmt := database.Builder.
		Delete("projects").
		Where("project_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("project_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("project_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to convert purge project query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}

	return nil
}

func (s *ProjectStore) applyQueryFilter(stmt sq.SelectBuilder, filter *types.ProjectFilter) sq.SelectBuilder {
	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("project_name", filter.Query))
	}
	//nolint:gocritic
	if filter.DeletedAt != nil {
		stmt = stmt.Where("project_deleted = ?", filter.DeletedAt)
	} else if filter.DeletedBeforeOrAt != nil {
		stmt = stmt.Where("project_deleted <= ?", filter.DeletedBeforeOrAt)
	} else {
		stmt = stmt.Where("project_deleted IS NULL")
	}

	if filter.TenantID != nil {
		stmt = stmt.Where("project_tenant_id = ?", filter.TenantID)
	}

	return stmt
}

func (s *ProjectStore) applySortFilter(stmt sq.SelectBuilder, filter *types.ProjectFilter) sq.SelectBuilder {
	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	switch filter.Sort {
	case enum.ProjectAttrCreated:
		stmt = stmt.OrderBy("project_created " + filter.Order.String())
	case enum.ProjectAttrUpdated:
		stmt = stmt.OrderBy("project_updated " + filter.Order.String())
	case enum.ProjectAttrDeleted:
		stmt = stmt.OrderBy("project_deleted " + filter.Order.String())
	}

	return stmt
}

func (s *ProjectStore) mapDBProjects(dbProjects []*project) []*types.Project {
	res := make([]*types.Project, len(dbProjects))
	for i := range dbProjects {
		res[i] = s.mapDBProject(dbProjects[i])
	}
	return res
}

func (s *ProjectStore) mapDBProject(dbProject *project) *types.Project {
	return &dbProject.Project
}
