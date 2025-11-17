package dbtx

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// ctxKeyTx is context key for storing and retrieving TransactionAccessor to and from a context.
type ctxKeyTx struct{}

// GetAccessor returns Accessor interface from the context if it exists or creates a new one from the provided *sql.DB.
// It is intended to be used in data layer functions that might or might not be running inside a transaction.
func GetAccessor(ctx context.Context, db *sqlx.DB) Accessor {
	if a, ok := ctx.Value(ctxKeyTx{}).(Accessor); ok {
		return a
	}
	return New(db)
}

// GetTransaction returns Transaction interface from the context if it exists or return nil.
// It is intended to be used in transactions in service layer functions to explicitly commit or rollback transactions.
func GetTransaction(ctx context.Context) Transaction {
	if a, ok := ctx.Value(ctxKeyTx{}).(Transaction); ok {
		return a
	}
	return nil
}
