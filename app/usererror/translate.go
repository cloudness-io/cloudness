package usererror

import (
	"context"
	"net/http"

	"github.com/cloudness-io/cloudness/blob"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types/check"

	"github.com/rs/zerolog/log"
)

func TranslateErrMsg(ctx context.Context, err error) string {
	return Translate(ctx, err).Error()
}

func Translate(ctx context.Context, err error) *Error {
	var (
		rError     *Error
		appError   *errors.Error
		checkError *check.ValidationError
		lockError  *lock.Error
	)

	// print original error for debugging purposes
	log.Ctx(ctx).Info().Err(err).Msgf("translating error to user facing error")

	// TODO: Improve performance of checking multiple
	// errors with errors.Is

	switch {
	// api errors
	case errors.As(err, &rError):
		return rError

	// status errors
	case errors.As(err, &appError):
		if appError.Err != nil {
			log.Ctx(ctx).Warn().Err(appError.Err).Msgf("Application error translation is omitting internal details.")
		}

		return &Error{Status: httpStatusCode(appError.Status), Message: appError.Message}

	// validation errors
	case errors.As(err, &checkError):
		return New(http.StatusBadRequest, checkError.Error())

	// store errors
	case errors.Is(err, store.ErrResourceNotFound):
		return ErrNotFound
	case errors.Is(err, store.ErrDuplicate):
		return ErrDuplicate
	case errors.Is(err, store.ErrPrimaryPathCantBeDeleted):
		return ErrPrimaryPathCantBeDeleted
	case errors.Is(err, store.ErrPathTooLong):
		return ErrPathTooLong
	case errors.Is(err, store.ErrNoChangeInRequestedMove):
		return ErrNoChange
	case errors.Is(err, store.ErrIllegalMoveCyclicHierarchy):
		return ErrCyclicHierarchy
	case errors.Is(err, store.ErrSpaceWithChildsCantBeDeleted):
		return ErrSpaceWithChildsCantBeDeleted

		//	upload
		//	errors
	case errors.Is(err, blob.ErrNotFound):
		return ErrNotFound

		//loc errors
	case errors.As(err, &lockError):
		return errorFromLockError(lockError)

		// unknown
		// error
	default:
		log.Ctx(ctx).Warn().Err(err).Msgf("Unable to translate error - returning Internal Error.")
		return ErrInternal
	}
}

// errorFromLockError returns the associated error for a given lock error.
func errorFromLockError(err *lock.Error) *Error {
	if err.Kind == lock.ErrorKindCannotLock ||
		err.Kind == lock.ErrorKindLockHeld ||
		err.Kind == lock.ErrorKindMaxRetriesExceeded {
		return ErrResourceLocked
	}

	return ErrInternal
}

var codes = map[errors.Status]int{
	errors.StatusConflict:           http.StatusConflict,
	errors.StatusInvalidArgument:    http.StatusBadRequest,
	errors.StatusNotFound:           http.StatusNotFound,
	errors.StatusNotImplemented:     http.StatusNotImplemented,
	errors.StatusPreconditionFailed: http.StatusPreconditionFailed,
	errors.StatusUnauthorized:       http.StatusUnauthorized,
	errors.StatusInternal:           http.StatusInternalServerError,
}

// httpStatusCode returns the associated HTTP status code for a status rror code.
func httpStatusCode(code errors.Status) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
