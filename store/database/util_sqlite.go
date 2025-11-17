//go:build !nosqlite
// +build !nosqlite

package database

import (
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func isSQLUniqueConstraintError(original error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(original, &sqliteErr) {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) ||
			errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey)
	}

	var pqErr *pq.Error
	if errors.As(original, &pqErr) {
		return pqErr.Code == "23505" // unique_violation
	}

	return false
}
