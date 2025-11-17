package database

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
)

var _ store.PrincipalStore = (*PrincipalStore)(nil)

// NewPrincipalStore returns a new PrincipalStore.
func NewPrincipalStore(db *sqlx.DB) *PrincipalStore {
	return &PrincipalStore{
		db: db,
	}
}

// PrincipalStore implements a PrincipalStore backed by a relational database.
type PrincipalStore struct {
	db *sqlx.DB
}

// principal is a DB representation of a principal.
// It is required to allow storing transformed UIDs used for uniquness constraints and searching.
type principal struct {
	types.Principal
}

// pricipalCommonColumns defines the columns that are the same across all principals.
const principalCommonColumns = `
	principal_id
	,principal_uid
	,principal_email
	,principal_display_name
	,principal_blocked
	,principal_salt
	,principal_created
	,principal_updated`

// principalColumns defines the column that are used only in a principal itself
// (for explicit principals the type is implicit, only the generic principal struct stores it explicitly).
const principalColumns = principalCommonColumns + `
,principal_type`

//nolint:goconst
const principalSelectBase = `
SELECT` + principalColumns + `
FROM principals`

// Find finds the principal by id.
func (s *PrincipalStore) Find(ctx context.Context, id int64) (*types.Principal, error) {
	const sqlQuery = principalSelectBase + `
		WHERE principal_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(principal)
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBPrincipal(dst), nil
}

func (s *PrincipalStore) mapDBPrincipal(dbPrincipal *principal) *types.Principal {
	return &dbPrincipal.Principal
}
