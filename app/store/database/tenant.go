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

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.TenantStore = (*TenantStore)(nil)

func NewTenantStore(db *sqlx.DB) *TenantStore {
	return &TenantStore{
		db: db,
	}
}

type TenantStore struct {
	db *sqlx.DB
}

type tenant struct {
	types.Tenant
}

const tenantInsert = `
INSERT INTO tenants (
	tenant_uid
	,tenant_name
	,tenant_slug
	,tenant_description
	,tenant_allow_admin_to_modify
	,tenant_max_projects
	,tenant_max_apps
	,tenant_max_instances
	,tenant_max_cpu_per_app
	,tenant_max_memory_per_app
	,tenant_max_volumes
	,tenant_min_volume_size
	,tenant_max_volume_size
	,tenant_created_by
	,tenant_created
	,tenant_updated
	,tenant_deleted
) values (
	:tenant_uid
	,:tenant_name
	,:tenant_slug
	,:tenant_description
	,:tenant_allow_admin_to_modify
	,:tenant_max_projects
	,:tenant_max_apps
	,:tenant_max_instances
	,:tenant_max_cpu_per_app
	,:tenant_max_memory_per_app
	,:tenant_max_volumes
	,:tenant_min_volume_size
	,:tenant_max_volume_size
	,:tenant_created_by
	,:tenant_created
	,:tenant_updated
	,:tenant_deleted
	) RETURNING tenant_id`

const tenantColumns = `
	tenant_id
	,tenant_uid
	,tenant_name
	,tenant_slug
	,tenant_description
	,tenant_allow_admin_to_modify
	,tenant_max_projects
	,tenant_max_apps
	,tenant_max_instances
	,tenant_max_cpu_per_app
	,tenant_max_memory_per_app
	,tenant_max_volumes
	,tenant_min_volume_size
	,tenant_max_volume_size
	,tenant_created_by
	,tenant_created
	,tenant_updated
	,tenant_deleted`

const tenantUpdateStmt = `Update tenants 
	SET 
		tenant_name = :tenant_name
		,tenant_description = :tenant_description
		,tenant_allow_admin_to_modify = :tenant_allow_admin_to_modify
		,tenant_max_projects = :tenant_max_projects
		,tenant_max_apps = :tenant_max_apps
		,tenant_max_instances = :tenant_max_instances
		,tenant_max_cpu_per_app = :tenant_max_cpu_per_app
		,tenant_max_memory_per_app = :tenant_max_memory_per_app
		,tenant_max_volumes = :tenant_max_volumes
		,tenant_min_volume_size = :tenant_min_volume_size
		,tenant_max_volume_size = :tenant_max_volume_size
		,tenant_updated = :tenant_updated
		,tenant_deleted = :tenant_deleted
	WHERE tenant_id = :tenant_id
`

func (s *TenantStore) Find(ctx context.Context, id int64) (*types.Tenant, error) {
	return s.FindDeleted(ctx, id, nil)
}

func (s *TenantStore) FindDeleted(ctx context.Context, id int64, deletedAt *int64) (*types.Tenant, error) {
	stmt := database.Builder.
		Select(tenantColumns).
		From("tenants").
		Where("tenant_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("tenant_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("tenant_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(tenant)
	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBTenant(dst), nil
}

func (s *TenantStore) FindByUID(ctx context.Context, tenantUID int64) (*types.Tenant, error) {
	return s.FindDeletedByUID(ctx, tenantUID, nil)
}

func (s *TenantStore) FindDeletedByUID(ctx context.Context, tenantUID int64, deletedAt **int64) (*types.Tenant, error) {
	stmt := database.Builder.
		Select(tenantColumns).
		From("tenants").
		Where("tenant_uid = ?", tenantUID)

	if deletedAt != nil {
		stmt = stmt.Where("tenant_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("tenant_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(tenant)
	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBTenant(dst), nil
}
func (s *TenantStore) Update(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	tenant.Updated = time.Now().UTC().UnixMilli()

	query, args, err := db.BindNamed(tenantUpdateStmt, tenant)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind tenant object")
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update query failed")
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return tenant, nil
}

func (s *TenantStore) Create(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(tenantInsert, tenant)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind tenant object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&tenant.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert query failed")
	}

	return tenant, nil
}

func (s *TenantStore) List(ctx context.Context, filter *types.TenantFilter) ([]*types.Tenant, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*tenant{}

	stmt := database.Builder.
		Select(tenantColumns).
		From("tenants")

	stmt = s.applyQueryFilter(stmt, filter)
	stmt = s.applySortFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by parent id query failed")

	}
	return s.mapDBTenants(dst), nil
}

// SoftDelete deletes the tenant softly by setting the deleted timestamp
func (s *TenantStore) SoftDelete(ctx context.Context, tenant *types.Tenant, deletedAt int64) error {
	tenant.Deleted = &deletedAt

	_, err := s.Update(ctx, tenant)
	return err
}

// Purge deletes the tenant permanently.
func (s *TenantStore) Purge(ctx context.Context, id int64, deletedAt *int64) error {
	stmt := database.Builder.
		Delete("tenants").
		Where("tenant_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("tenant_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("tenant_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to convert purge tenant query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}

	return nil
}

func (s *TenantStore) applyQueryFilter(stmt squirrel.SelectBuilder, filter *types.TenantFilter) squirrel.SelectBuilder {
	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("tenant_name", filter.Query))
	}
	//nolint:gocritic
	if filter.DeletedAt != nil {
		stmt = stmt.Where("tenant_deleted = ?", filter.DeletedAt)
	} else if filter.DeletedBeforeOrAt != nil {
		stmt = stmt.Where("tenant_deleted <= ?", filter.DeletedBeforeOrAt)
	} else {
		stmt = stmt.Where("tenant_deleted IS NULL")
	}

	return stmt
}

func (s *TenantStore) applySortFilter(stmt squirrel.SelectBuilder, filter *types.TenantFilter) squirrel.SelectBuilder {
	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	switch filter.Sort {
	case enum.TenantAttrCreated:
		stmt = stmt.OrderBy("project_created " + filter.Order.String())
	case enum.TenantAttrUpdated:
		stmt = stmt.OrderBy("project_updated " + filter.Order.String())
	case enum.TenantAttrDeleted:
		stmt = stmt.OrderBy("project_deleted " + filter.Order.String())
	}

	return stmt
}

func (s *TenantStore) mapDBTenants(dbTenants []*tenant) []*types.Tenant {
	res := make([]*types.Tenant, len(dbTenants))
	for i := range dbTenants {
		res[i] = s.mapDBTenant(dbTenants[i])
	}
	return res
}

func (s *TenantStore) mapDBTenant(dbTenant *tenant) *types.Tenant {
	return &dbTenant.Tenant
}
