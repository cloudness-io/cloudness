package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloudness-io/cloudness/store"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// default query range limit.
const defaultLimit = 100

// limit returns the page size to a sql limit.
func Limit(size int) uint64 {
	if size == 0 {
		size = defaultLimit
	}
	return uint64(size)
}

// offset converts the page to a sql offset.
func Offset(page, size int) uint64 {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = defaultLimit
	}
	page--
	return uint64(page * size)
}

// Logs the error and message, returns either the provided message.
// Always logs the full message with error as warning.
//
//nolint:unparam // revisit error processing
func ProcessSQLErrorf(ctx context.Context, err error, format string, args ...interface{}) error {
	// If it's a known error, return converted error instead.
	translatedError := err
	switch {
	case errors.Is(err, sql.ErrNoRows):
		translatedError = store.ErrResourceNotFound
	case isSQLUniqueConstraintError(err):
		translatedError = store.ErrDuplicate
	default:
	}

	//nolint:errorlint // we want to match exactly here.
	if translatedError != err {
		log.Ctx(ctx).Debug().Err(err).Msgf("translated sql error to: %s", translatedError)
	}

	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), translatedError)
}
