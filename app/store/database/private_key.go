package database

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
)

var _ store.PrivateKeyStore = (*PrivateKeyStore)(nil)

func NewPrivateKeyStore(db *sqlx.DB) *PrivateKeyStore {
	return &PrivateKeyStore{
		db: db,
	}
}

type PrivateKeyStore struct {
	db *sqlx.DB
}

type privateKey struct {
	types.PrivateKey
}

const privateKeyColumns = `
	private_key_id,
	private_key_uid,
	private_key_tenant_id,
	private_key_name,
	private_key_description,
	private_key_pem,
	private_key_is_git,
	private_key_created,
	private_key_updated`

const privateKeyInsert = `
INSERT INTO private_keys (
	private_key_uid
	,private_key_tenant_id
	,private_key_name
	,private_key_description
	,private_key_pem
	,private_key_is_git
	,private_key_created
	,private_key_updated	
) values (
	:private_key_uid
	,:private_key_tenant_id
	,:private_key_name
	,:private_key_description
	,:private_key_pem
	,:private_key_is_git
	,:private_key_created
	,:private_key_updated
	) RETURNING private_key_id
	`

const privateKeySelectBase = `
	SELECT` + privateKeyColumns + `
	FROM private_keys`

// Find the private key by id.
func (s *PrivateKeyStore) Find(ctx context.Context, tenantID, id int64) (*types.PrivateKey, error) {
	const sqlQuery = privateKeySelectBase + `
	WHERE private_key_tenant_id = $1 AND private_key_id = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(privateKey)

	if err := db.GetContext(ctx, dst, sqlQuery, tenantID, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select private key by id query failed")
	}
	return s.mapDBPrivateKey(dst), nil
}

// Create save the private key
func (s *PrivateKeyStore) Create(ctx context.Context, privatekey *types.PrivateKey) (*types.PrivateKey, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(privateKeyInsert, privatekey)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind private key object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&privatekey.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert private key query failed")
	}

	return privatekey, nil
}

// Delete delets the private key
func (s *PrivateKeyStore) Delete(ctx context.Context, tenantID, id int64) error {
	privateKeyDeleteStmt := `DELETE FROM private_keys WHERE private_key_tenant_id = $1 AND private_key_id = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, privateKeyDeleteStmt, tenantID, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete private key query failed")
	}
	return nil
}

func (s *PrivateKeyStore) mapDBPrivateKey(dst *privateKey) *types.PrivateKey {
	return &dst.PrivateKey
}
