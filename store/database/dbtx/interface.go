package dbtx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Accessor is the SQLx database manipulation interface.
type Accessor interface {
	sqlx.ExtContext // sqlx.binder + sqlx.QueryerContext + sqlx.ExecerContext
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Transaction is the Go's standard sql transaction interface.
type Transaction interface {
	Commit() error
	Rollback() error
}

type Transactor interface {
	WithTx(ctx context.Context, txFn func(ctx context.Context) error, opts ...interface{}) error
}

// AccessorTx is used to access the database. It combines Accessor interface
// with Transactor (capability to run functions in a transaction).
type AccessorTx interface {
	Accessor
	Transactor
}

// TransactionAccessor combines data access capabilities with the transaction commit and rollback.
type TransactionAccessor interface {
	Transaction
	Accessor
}
