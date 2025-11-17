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

var _ store.VolumeStore = (*VolumeStore)(nil)

// NewVolumeStore returns new VolumeStore
func NewVolumeStore(db *sqlx.DB) *VolumeStore {
	return &VolumeStore{db: db}
}

// VolumeStore implements a VolumeStore backed by a relational database.
type VolumeStore struct {
	db *sqlx.DB
}

// volume is a DB representation of a volume
type volume struct {
	types.Volume
}

// volumeColumns defines the columns of the volumes table.
const volumeColumns = `
	 volume_id
	,volume_uid
	,volume_tenant_id
   	,volume_project_id
   	,volume_environment_id
	,volume_environment_uid
	,volume_server_id
   	,volume_application_id
   	,volume_name
   	,volume_mount_path
   	,volume_host_path
   	,volume_size
   	,volume_created
   	,volume_updated
   	,volume_deleted
	`

const volumeInsert = `
INSERT INTO volumes(
	volume_uid
	,volume_tenant_id
	,volume_project_id
	,volume_environment_id
	,volume_environment_uid
	,volume_server_id
	,volume_application_id
	,volume_name 
	,volume_mount_path
	,volume_host_path
	,volume_size
	,volume_created
	,volume_updated
	,volume_deleted
) VALUES (
	:volume_uid
	,:volume_tenant_id
	,:volume_project_id
	,:volume_environment_id
	,:volume_environment_uid
	,:volume_server_id
	,:volume_application_id
	,:volume_name
	,:volume_mount_path
	,:volume_host_path
	,:volume_size
	,:volume_created
	,:volume_updated
	,:volume_deleted
) RETURNING volume_id`

const volumeUpdate = `UPDATE volumes 
	SET 
		volume_name = :volume_name
		,volume_application_id = :volume_application_id
		,volume_mount_path = :volume_mount_path
		,volume_host_path = :volume_host_path
		,volume_size = :volume_size
		,volume_updated = :volume_updated
		,volume_deleted = :volume_deleted
	WHERE volume_id = :volume_id AND volume_project_id = :volume_project_id AND volume_tenant_id = :volume_tenant_id AND volume_environment_id = :volume_environment_id`

// Create creates a new volume
func (s *VolumeStore) Create(ctx context.Context, volume *types.Volume) (*types.Volume, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	query, arg, err := db.BindNamed(volumeInsert, volume)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind volume object")
	}
	if err = db.QueryRowContext(ctx, query, arg...).Scan(&volume.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert volume query failed")
	}
	return volume, nil
}

// Update updates the volume
func (s *VolumeStore) Update(ctx context.Context, volume *types.Volume) (*types.Volume, error) {
	volume.Updated = time.Now().UTC().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(volumeUpdate, volume)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind volume object")
	}
	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update volume query failed")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return volume, nil
}

// Find the volume by id
func (s *VolumeStore) Find(ctx context.Context, id int64) (*types.Volume, error) {
	return s.FindDeleted(ctx, id, nil)

}
func (s *VolumeStore) FindDeleted(ctx context.Context, id int64, deletedAt *int64) (*types.Volume, error) {
	stms := database.Builder.
		Select(volumeColumns).
		From("volumes").
		Where("volume_id = ?", id)

	if deletedAt != nil {
		stms = stms.Where("volume_deleted = ?", deletedAt)
	} else {
		stms = stms.Where("volume_deleted IS NULL")
	}
	sql, args, err := stms.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(volume)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDBVolume(dst)
}

// FindByUID finds the volume by tenant id, project id ,environment id, volume id and volume u_id
func (s *VolumeStore) FindByUID(ctx context.Context, tenantID, projectID, environmentID, uid int64) (*types.Volume, error) {
	return s.FindDeletedByUID(ctx, tenantID, projectID, environmentID, uid, nil)
}

func (s *VolumeStore) FindDeletedByUID(ctx context.Context, tenantID, projectID, environmentID, uid int64, deletedAt *int64) (*types.Volume, error) {
	stmt := database.Builder.
		Select(volumeColumns).
		From("volumes").
		Where("volume_tenant_id = ?", tenantID).
		Where("volume_project_id = ?", projectID).
		Where("volume_environment_id = ?", environmentID).
		Where("volume_uid = ?", uid)

	if deletedAt != nil {
		stmt = stmt.Where("volume_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("volume_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(volume)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by uid query failed")
	}
	return s.mapDBVolume(dst)
}

// List lists the volumes by filter
func (s *VolumeStore) List(ctx context.Context, filter *types.VolumeFilter) ([]*types.Volume, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*volume{}

	stmt := database.Builder.
		Select(volumeColumns).
		From("volumes")

	stmt = s.applyQueryFilter(stmt, filter)
	stmt = s.applySortFilter(stmt, filter)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List by parent id query failed")

	}
	return s.mapDBVolumes(dst)

}

func (s *VolumeStore) SoftDelete(ctx context.Context, volume *types.Volume, deletedAt int64) error {
	volume.Deleted = &deletedAt

	_, err := s.Update(ctx, volume)
	return err
}

func (s *VolumeStore) Purge(ctx context.Context, id int64, deletedAt *int64) error {
	stmt := database.Builder.
		Delete("volumes").
		Where("volume_id = ?", id)

	if deletedAt != nil {
		stmt = stmt.Where("volume_deleted = ?", *deletedAt)
	} else {
		stmt = stmt.Where("volume_deleted IS NULL")
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to convert purge volume query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}

	return nil
}
func (s *VolumeStore) applyQueryFilter(stmt squirrel.SelectBuilder, filter *types.VolumeFilter) squirrel.SelectBuilder {
	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("volume_name", filter.Query))
	}

	//nolint:gocritic
	if filter.DeletedAt != nil {
		stmt = stmt.Where("volume_deleted = ?", filter.DeletedAt)
	} else if filter.DeletedBeforeOrAt != nil {
		stmt = stmt.Where("volume_deleted <= ?", filter.DeletedBeforeOrAt)
	} else {
		stmt = stmt.Where("volume_deleted IS NULL")
	}

	if filter.TenantID != nil {
		stmt = stmt.Where("volume_tenant_id = ?", filter.TenantID)
	}

	if filter.ProjectID != nil {
		stmt = stmt.Where("volume_project_id = ?", filter.ProjectID)
	}

	if filter.EnvironmentID != nil {
		stmt = stmt.Where("volume_environment_id = ?", filter.EnvironmentID)
	}

	if filter.ApplicationID != nil {
		stmt = stmt.Where("volume_application_id = ?", filter.ApplicationID)
	} else {
		stmt = stmt.Where("volume_application_id IS NULL")
	}

	if filter.ServerID != nil {
		stmt = stmt.Where("volume_server_id = ?", filter.ServerID)
	}

	return stmt
}

func (s *VolumeStore) applySortFilter(stmt squirrel.SelectBuilder, filter *types.VolumeFilter) squirrel.SelectBuilder {
	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	switch filter.Sort {
	case enum.VolumeAttrNone:
		stmt = stmt.OrderBy("volume_id " + enum.OrderAsc.String())
	case enum.VolumeAttrCreated:
		stmt = stmt.OrderBy("volume_created " + filter.Order.String())
	case enum.VolumeAttrUpdated:
		stmt = stmt.OrderBy("volume_updated" + filter.Order.String())
	case enum.VolumeAttrDeleted:
		stmt = stmt.OrderBy("volume_deleted" + filter.Order.String())
	}

	return stmt
}

func (s *VolumeStore) mapDBVolumes(dbVolumes []*volume) ([]*types.Volume, error) {
	res := make([]*types.Volume, len(dbVolumes))
	for i := range dbVolumes {
		c, err := s.mapDBVolume(dbVolumes[i])
		if err != nil {
			return nil, err
		}
		res[i] = c
	}
	return res, nil
}

func (s *VolumeStore) mapDBVolume(dbVolume *volume) (*types.Volume, error) {
	return &dbVolume.Volume, nil
}
