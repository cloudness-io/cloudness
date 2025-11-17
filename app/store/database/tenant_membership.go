package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.TenantMembershipStore = (*TenantMembershipStore)(nil)

type TenantMembershipStore struct {
	db *sqlx.DB
}

func NewTenantMembershipStore(db *sqlx.DB) *TenantMembershipStore {
	return &TenantMembershipStore{
		db: db,
	}
}

type tenantMembershipDB struct {
	ID          int64           `db:"tenant_membership_id"`
	TenantID    int64           `db:"tenant_membership_tenant_id"`
	PrincipalID int64           `db:"tenant_membership_principal_id"`
	Role        enum.TenantRole `db:"tenant_membership_role"`
	CreatedBy   int64           `db:"tenant_membership_created_by"`
	Created     int64           `db:"tenant_membership_created"`
	Updated     int64           `db:"tenant_membership_updated"`
}

type tenantMembership struct {
	tenantMembershipDB
	types.Principal
	types.Tenant
}

const tenantMembershipColumns = `
	tenant_membership_id,
	tenant_membership_tenant_id,
	tenant_membership_principal_id,
	tenant_membership_role,
	tenant_membership_created_by,
	tenant_membership_created,
	tenant_membership_updated
`

const tenantUserColums = `
	principal_email,
	principal_display_name 
`

const tenantMembershipSelect = `
	SELECT ` + tenantMembershipColumns + `
	FROM tenant_memberships`

const tenantMembershipInsert = `
INSERT INTO tenant_memberships (
	tenant_membership_tenant_id
	,tenant_membership_principal_id
	,tenant_membership_role
	,tenant_membership_created_by
	,tenant_membership_created
	,tenant_membership_updated
) VALUES (
	:tenant_membership_tenant_id
	,:tenant_membership_principal_id
	,:tenant_membership_role
	,:tenant_membership_created_by
	,:tenant_membership_created
	,:tenant_membership_updated
	)`

func mapToInternalTenantMembership(m *types.TenantMembership) tenantMembershipDB {
	return tenantMembershipDB{
		TenantID:    m.TenantID,
		PrincipalID: m.PrincipalID,
		Role:        m.Role,
		CreatedBy:   m.CreatedBy,
		Created:     m.Created,
		Updated:     m.Updated,
	}
}

func (s *TenantMembershipStore) Find(ctx context.Context, tenantID, principalID int64) (*types.TenantMembership, error) {
	stmt := database.Builder.Select(tenantMembershipColumns).
		From("tenant_memberships").
		Where("tenant_membership_tenant_id = ?", tenantID).
		Where("tenant_membership_principal_id = ?", principalID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(tenantMembership)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select membership by tenant and principal id query failed")
	}
	return s.mapToMembership(dst), nil
}

func (s *TenantMembershipStore) List(ctx context.Context, principalID int64) ([]*types.TenantMembershipUser, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*tenantMembership{}
	const columns = tenantMembershipColumns + "," + tenantColumns + "," + tenantUserColums
	stmt := database.Builder.
		Select(columns).
		From("tenant_memberships").
		InnerJoin("tenants ON tenants.tenant_id = tenant_membership_tenant_id").
		InnerJoin("principals ON principals.principal_id = tenant_membership_principal_id").
		Where("tenant_membership_principal_id = $1")

	sql, _, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, principalID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by principal id query failed")
	}

	return s.mapToMembershipTenants(dst), nil
}

func (s *TenantMembershipStore) ListByTenant(ctx context.Context, principalID int64) ([]*types.TenantMembershipUser, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*tenantMembership{}

	const columns = tenantMembershipColumns + "," + tenantColumns + "," + tenantUserColums
	stmt := database.Builder.
		Select(columns).
		From("tenant_memberships").
		InnerJoin("tenants ON tenants.tenant_id = tenant_membership_tenant_id").
		InnerJoin("principals ON principals.principal_id = tenant_membership_principal_id").
		Where("tenant_membership_tenant_id = $1")

	sql, _, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, principalID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by tenant id query failed")
	}

	return s.mapToMembershipTenants(dst), nil
}

func (s *TenantMembershipStore) Create(ctx context.Context, tenantMembership *types.TenantMembership) error {
	now := time.Now().UTC().UnixMilli()
	tenantMembership.Created = now
	tenantMembership.Updated = now

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(tenantMembershipInsert, mapToInternalTenantMembership(tenantMembership))
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind tenent membership object")
	}
	if _, err = db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert query failed")
	}
	return nil
}

func (s *TenantMembershipStore) Update(ctx context.Context, tenantID, principalID int64, role enum.TenantRole) error {
	stmt := database.Builder.Update("tenant_memberships").
		Where("tenant_membership_tenant_id = ?", tenantID).
		Where("tenant_membership_principal_id = ?", principalID).
		Set("tenant_membership_role", string(role))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the update query failed")
	}
	return nil
}

func (s *TenantMembershipStore) Delete(ctx context.Context, tenantID, principalID int64) error {
	stmt := database.Builder.Delete("tenant_memberships").
		Where("tenant_membership_tenant_id = ?", tenantID).
		Where("tenant_membership_principal_id = ?", principalID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}
	return nil
}

func (s *TenantMembershipStore) DeleteAll(ctx context.Context, tenantID int64) error {
	stmt := database.Builder.Delete("tenant_memberships").
		Where("tenant_membership_tenant_id = ?", tenantID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}
	return nil
}

func (s *TenantMembershipStore) mapToMembershipTenants(mts []*tenantMembership) []*types.TenantMembershipUser {
	res := make([]*types.TenantMembershipUser, len(mts))

	for i := range mts {
		t := mts[i]
		res[i] = s.mapToMembershipTenant(t)
	}

	return res
}

func (s *TenantMembershipStore) mapToMembershipTenant(mt *tenantMembership) *types.TenantMembershipUser {
	return &types.TenantMembershipUser{
		Email:            mt.Email,
		DisplayName:      mt.DisplayName,
		TenantMembership: s.mapToMembership(mt),
	}
}

func (s *TenantMembershipStore) mapToMembership(mt *tenantMembership) *types.TenantMembership {
	return &types.TenantMembership{
		TenantMembershipKey: &types.TenantMembershipKey{
			ID:                mt.tenantMembershipDB.ID,
			TenantUID:         mt.Tenant.UID,
			TenantID:          mt.TenantID,
			TenantName:        mt.Tenant.Name,
			TenantDescription: mt.Tenant.Description,
			PrincipalID:       mt.PrincipalID,
			Role:              mt.Role,
		},
		Created: mt.tenantMembershipDB.Created,
		Updated: mt.tenantMembershipDB.Updated,
	}
}
