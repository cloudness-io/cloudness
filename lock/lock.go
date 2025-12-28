package lock

import (
	"context"
	"fmt"
)

// ErrorKind enum displays human readable message
// in error.
type ErrorKind string

const (
	ErrorKindLockHeld            ErrorKind = "lock already held"
	ErrorKindLockNotHeld         ErrorKind = "lock not held"
	ErrorKindProviderError       ErrorKind = "lock provider error"
	ErrorKindCannotLock          ErrorKind = "timeout while trying to acquire lock"
	ErrorKindContext             ErrorKind = "context error while trying to acquire lock"
	ErrorKindMaxRetriesExceeded  ErrorKind = "max retries exceeded to acquire lock"
	ErrorKindGenerateTokenFailed ErrorKind = "token generation failed"
)

// Error is custom unique type for all type of errors.
type Error struct {
	Kind ErrorKind
	Key  string
	Err  error
}

func NewError(kind ErrorKind, key string, err error) *Error {
	return &Error{
		Kind: kind,
		Key:  key,
		Err:  err,
	}
}

// Error implements error interface.
func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s on key %s with err: %v", e.Kind, e.Key, e.Err)
	}
	return fmt.Sprintf("%s on key %s", e.Kind, e.Key)
}

// MutexManager describes a Distributed Lock Manager.
type MutexManager interface {
	// NewMutex creates a mutex for the given key. The returned mutex is not held
	// and must be acquired with a call to .Lock.
	NewMutex(key string, options ...Option) (Mutex, error)
}

type Mutex interface {
	// Key returns the key to be locked.
	Key() string

	// Lock acquires the lock. It fails with error if the lock is already held.
	Lock(ctx context.Context) error

	// Unlock releases the lock. It fails with error if the lock is not currently held.
	Unlock(ctx context.Context) error

	// IsHeld returns true if the lock is currently held.
	IsHeld(ctx context.Context) bool
}
