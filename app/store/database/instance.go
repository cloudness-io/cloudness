package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.InstanceStore = (*InstanceStore)(nil)

// NewInstanceStore creates a new instance store.
func NewInstanceStore(db *sqlx.DB) *InstanceStore {
	return &InstanceStore{db: db}
}

// InstanceStore implements a InstanceStore backed by a relational database.
type InstanceStore struct {
	db *sqlx.DB
}

// instance is a DB representation of a instance
type instance struct {
	types.Instance
}

const instanceColumns = `
	instance_id
	,instance_super_admin
	,instance_public_ipv4
	,instance_public_ipv6
	,instance_fqdn
	,instance_update_enabled
	,instance_update_check_frequency
	,instance_dns_validation_enabled
	,instance_dns_servers
	,instance_user_signup_enabled
	,instance_demo_user_enabled
	,instance_registry_enabled
	,instance_registry_size
	,instance_registry_mirror_enabled
	,instance_registry_mirror_size
	,instance_external_scripts
	,instance_created
	,instance_updated`

const instanceInsert = `
INSERT INTO instances (
	instance_public_ipv4
	,instance_public_ipv6
	,instance_fqdn
	,instance_update_enabled
	,instance_update_check_frequency
	,instance_dns_validation_enabled
	,instance_dns_servers
	,instance_user_signup_enabled
	,instance_demo_user_enabled
	,instance_registry_enabled
	,instance_registry_size
	,instance_registry_mirror_enabled
	,instance_registry_mirror_size
	,instance_external_scripts
	,instance_created
	,instance_updated
) VALUES (
	:instance_public_ipv4
	,:instance_public_ipv6
	,:instance_fqdn
	,:instance_update_enabled
	,:instance_update_check_frequency
	,:instance_dns_validation_enabled
	,:instance_dns_servers
	,:instance_user_signup_enabled
	,:instance_demo_user_enabled
	,:instance_registry_enabled
	,:instance_registry_size
	,:instance_registry_mirror_enabled
	,:instance_registry_mirror_size
	,:instance_external_scripts
	,:instance_created
	,:instance_updated
) RETURNING instance_id
`

// Create saves the instance settings.
func (s *InstanceStore) Create(ctx context.Context, instance *types.Instance) (*types.Instance, error) {
	now := time.Now().UTC().UnixMilli()
	instance.Created = now
	instance.Updated = now
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(instanceInsert, instance)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind instance oject")
	}

	if err = db.QueryRowContext(ctx, query, args...).Scan(&instance.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert instance query failed")
	}
	return instance, nil
}

// Get gets the only instance settings
func (s *InstanceStore) Get(ctx context.Context) (*types.Instance, error) {
	stmt := database.Builder.
		Select(instanceColumns).
		From("instances").
		Where("instance_id = ?", 1)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(instance)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDbInstance(dst), nil
}

// Update updates the instance settings
func (s *InstanceStore) Update(ctx context.Context, instance *types.Instance) (*types.Instance, error) {
	const instanceUpdate = `UPDATE instances
 	SET
		instance_super_admin = :instance_super_admin
		,instance_public_ipv4 = :instance_public_ipv4
		,instance_public_ipv6 = :instance_public_ipv6
		,instance_fqdn = :instance_fqdn
		,instance_update_enabled = :instance_update_enabled
		,instance_update_check_frequency = :instance_update_check_frequency
		,instance_dns_validation_enabled = :instance_dns_validation_enabled
		,instance_dns_servers = :instance_dns_servers
		,instance_user_signup_enabled = :instance_user_signup_enabled
		,instance_demo_user_enabled = :instance_demo_user_enabled
		,instance_registry_enabled = :instance_registry_enabled
		,instance_registry_size = :instance_registry_size
		,instance_registry_mirror_enabled = :instance_registry_mirror_enabled
		,instance_registry_mirror_size = :instance_registry_mirror_size
		,instance_external_scripts = :instance_external_scripts
		,instance_created = :instance_created
		,instance_updated = :instance_updated
	WHERE instance_id = 1`

	db := dbtx.GetAccessor(ctx, s.db)
	query, args, err := db.BindNamed(instanceUpdate, instance)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind instance object")
	}

	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update instance query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}
	return instance, nil
}

func (s *InstanceStore) mapDbInstance(i *instance) *types.Instance {
	return &i.Instance
}
