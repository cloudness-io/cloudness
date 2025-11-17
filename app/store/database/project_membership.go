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

var _ store.ProjectMembershipStore = (*ProjectMembershipStore)(nil)

type ProjectMembershipStore struct {
	db *sqlx.DB
}

func NewProjectMembershipStore(db *sqlx.DB) *ProjectMembershipStore {
	return &ProjectMembershipStore{
		db: db,
	}
}

type projectMembershipDB struct {
	TenatID            int64            `db:"project_membership_tenant_id"`
	TenantMembershipID int64            `db:"project_membership_tenant_membership_id"`
	ProjectID          int64            `db:"project_membership_project_id"`
	PrincipalID        int64            `db:"project_membership_principal_id"`
	Role               enum.ProjectRole `db:"project_membership_role"`
	CreatedBy          int64            `db:"project_membership_created_by"`

	Created int64 `db:"project_membership_created"`
	Updated int64 `db:"project_membership_updated"`
}

type projectMembership struct {
	projectMembershipDB
	types.Principal
}

const projectMembershipColumns = `
	project_membership_tenant_id,
	project_membership_tenant_membership_id,
	project_membership_project_id,
	project_membership_principal_id,
	project_membership_role,
	project_membership_created_by,
	project_membership_created,
	project_membership_updated
`

const projectMembershipInsert = `
INSERT INTO project_memberships (
	project_membership_tenant_id
	,project_membership_tenant_membership_id
	,project_membership_project_id
	,project_membership_principal_id
	,project_membership_role
	,project_membership_created_by
	,project_membership_created
	,project_membership_updated
) VALUES (
	:project_membership_tenant_id
	,:project_membership_tenant_membership_id
	,:project_membership_project_id
	,:project_membership_principal_id
	,:project_membership_role
	,:project_membership_created_by
	,:project_membership_created
	,:project_membership_updated
	)`

func mapToInternalProjectMemerbship(m *types.ProjectMembership) projectMembershipDB {
	return projectMembershipDB{
		TenatID:            m.TenantID,
		TenantMembershipID: m.TenantMembershipID,
		ProjectID:          m.ProjectID,
		PrincipalID:        m.PrincipalID,
		Role:               m.Role,
		CreatedBy:          m.CreatedBy,
		Created:            m.Created,
		Updated:            m.Updated,
	}
}

func (s *ProjectMembershipStore) Find(ctx context.Context, tenantID, projectID, principalID int64) (*types.ProjectMembership, error) {
	stmt := database.Builder.Select(projectMembershipColumns).
		From("project_memberships").
		Where("project_membership_tenant_id = ?", tenantID).
		Where("project_membership_project_id = ?", projectID).
		Where("project_membership_principal_id = ?", principalID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(projectMembership)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select project membership by tenant project and principal id query failed")
	}
	return s.mapToProjectMembership(dst), nil

}

func (s *ProjectMembershipStore) List(ctx context.Context, tenantID int64, projectID int64) ([]*types.ProjectMembershipUser, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*projectMembership{}

	stmt := database.Builder.
		Select(projectMembershipColumns+principalColumns).
		From("project_memberships").
		InnerJoin("principals ON principals.principal_id = project_membership_principal_id").
		Where("project_membership_tenant_id = ?", tenantID).
		Where("project_membership_project_id = ?", projectID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List project membership query failed")
	}

	return s.mapToProjectUserMemberships(dst), nil
}

func (s *ProjectMembershipStore) Create(ctx context.Context, projectMembership *types.ProjectMembership) error {
	now := time.Now().UTC().UnixMilli()
	projectMembership.Created = now
	projectMembership.Updated = now

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(projectMembershipInsert, mapToInternalProjectMemerbship(projectMembership))
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind project membership object")
	}
	if _, err = db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert project membership query failed")
	}
	return nil
}

func (s *ProjectMembershipStore) Update(ctx context.Context, tenantID, projectID, principalID int64, role enum.ProjectRole) error {
	stmt := database.Builder.Update("project_memberships").
		Where("project_membership_tenant_id = ?", tenantID).
		Where("project_membership_project_id = ?", projectID).
		Where("project_membership_principal_id = ?", principalID).
		Set("project_membership_role", string(role))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update project membership query failed")
	}
	return nil
}

func (s *ProjectMembershipStore) Delete(ctx context.Context, tenantID, projectID, principalID int64) error {
	stmt := database.Builder.Delete("project_memberships").
		Where("project_membership_tenant_id = ?", tenantID).
		Where("project_membership_project_id = ?", projectID).
		Where("project_membership_principal_id = ?", principalID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete project membership query failed")
	}
	return nil
}

func (s *ProjectMembershipStore) mapToProjectUserMembership(pm *projectMembership) *types.ProjectMembershipUser {
	return &types.ProjectMembershipUser{
		ProjectMembership: s.mapToProjectMembership(pm),
		DisplayName:       pm.DisplayName,
		Email:             pm.Email,
	}
}

func (s *ProjectMembershipStore) mapToProjectUserMemberships(pms []*projectMembership) []*types.ProjectMembershipUser {
	res := make([]*types.ProjectMembershipUser, len(pms))
	for i := range pms {
		p := pms[i]
		res[i] = s.mapToProjectUserMembership(p)
	}
	return res
}

func (s *ProjectMembershipStore) mapToProjectMembership(pm *projectMembership) *types.ProjectMembership {
	return &types.ProjectMembership{
		ProjectMembershipKey: &types.ProjectMembershipKey{
			TenantID:           pm.TenatID,
			TenantMembershipID: pm.TenantMembershipID,
			ProjectID:          pm.ProjectID,
			PrincipalID:        pm.PrincipalID,
			Role:               pm.Role,
		},
		Created: pm.projectMembershipDB.Created,
		Updated: pm.projectMembershipDB.Updated,
	}
}
