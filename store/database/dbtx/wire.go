package dbtx

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideAccessorTx,
	ProvideAccessor,
	ProvideTransactor,
)

// ProvideAccessorTx provides the most versatile database access interface.
// All DB queries and transactions can be performed.
func ProvideAccessorTx(db *sqlx.DB) AccessorTx {
	return New(db)
}

// ProvideAccessor provides the database access interface. All DB queries can be performed.
func ProvideAccessor(a AccessorTx) Accessor {
	return a
}

// ProvideTransactor provides ability to run DB transactions.
func ProvideTransactor(a AccessorTx) Transactor {
	return a
}
