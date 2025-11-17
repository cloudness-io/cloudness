package dbtx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// New returns new database Runner interface.
func New(db *sqlx.DB) AccessorTx {
	mx := getLocker(db)
	run := &runnerDB{
		db: sqlDB{db},
		mx: mx,
	}
	return run
}

// transactor is combines data access capabilities with transaction starting.
type transactor interface {
	Accessor
	startTx(ctx context.Context, opts *sql.TxOptions) (TransactionAccessor, error)
}

// sqlDB is a wrapper for the sqlx.DB that implements the transactor interface.
type sqlDB struct {
	*sqlx.DB
}

var _ transactor = (*sqlDB)(nil)

func (db sqlDB) startTx(ctx context.Context, opts *sql.TxOptions) (TransactionAccessor, error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	return tx, err
}
