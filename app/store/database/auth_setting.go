package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.AuthSettingsStore = (*AuthSettingsStore)(nil)

type AuthSettingsStore struct {
	db *sqlx.DB
}

func NewAuthSettingsStore(db *sqlx.DB) *AuthSettingsStore {
	return &AuthSettingsStore{
		db: db,
	}
}

type authSetting struct {
	types.AuthSetting
}

const authSettingColumns = `
	auth_id,
	auth_provider,
   auth_enabled,
   auth_client_id,
   auth_client_secret,
   auth_base_url,
   auth_created,
   auth_updated
`

const authSettingsInsert = `
INSERT INTO auth_settings (
	auth_provider
	,auth_enabled
	,auth_client_id
	,auth_client_secret
	,auth_base_url
	,auth_created
	,auth_updated
) VALUES (
	:auth_provider
	,:auth_enabled
	,:auth_client_id
	,:auth_client_secret
	,:auth_base_url
	,:auth_created
	,:auth_updated
) 
ON CONFLICT (auth_provider) DO NOTHING
RETURNING auth_id`

// Create saves the auth setting.
func (s *AuthSettingsStore) Create(ctx context.Context, auth *types.AuthSetting) (*types.AuthSetting, error) {
	auth.Created = time.Now().UTC().UnixMilli()
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(authSettingsInsert, auth)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind auth setting object")
	}

	if _, err = db.ExecContext(ctx, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to create auth setting object")
	}
	return auth, nil
}

// Update updates the auth settings
func (s *AuthSettingsStore) Update(ctx context.Context, auth *types.AuthSetting) (*types.AuthSetting, error) {
	const updateStmt = `UPDATE auth_settings
	SET 
		auth_enabled = :auth_enabled
		,auth_client_id = :auth_client_id
		,auth_client_secret = :auth_client_secret
		,auth_base_url = :auth_base_url
		,auth_updated = :auth_updated
	WHERE auth_provider = :auth_provider`

	auth.Updated = time.Now().UTC().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)
	query, args, err := db.BindNamed(updateStmt, auth)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind auth setting object")
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to update auth setting object")
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return auth, nil
}

// FindByProvider gets the auth settings by provider.
func (s *AuthSettingsStore) FindByProvider(ctx context.Context, provider enum.AuthProvider) (*types.AuthSetting, error) {
	stmt := database.Builder.
		Select(authSettingColumns).
		From("auth_settings").
		Where("auth_provider = ?", provider)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(authSetting)
	if err := db.GetContext(ctx, dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by provider query failed")
	}

	return s.mapDBAuthSetting(dst), nil
}

// List lists the auth settings.
func (s *AuthSettingsStore) List(ctx context.Context) ([]*types.AuthSetting, error) {
	stmt := database.Builder.
		Select(authSettingColumns).
		From("auth_settings")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*authSetting{}
	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select list query failed")
	}

	return s.mapDBAuthSettings(dst), nil
}

func (s *AuthSettingsStore) mapDBAuthSetting(d *authSetting) *types.AuthSetting {
	return &d.AuthSetting
}

func (s *AuthSettingsStore) mapDBAuthSettings(dbAuthSettings []*authSetting) []*types.AuthSetting {
	dep := make([]*types.AuthSetting, len(dbAuthSettings))
	for i := range dbAuthSettings {
		dep[i] = s.mapDBAuthSetting(dbAuthSettings[i])
	}
	return dep
}
