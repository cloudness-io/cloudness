package check

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrAny  = &ValidationError{}
	ErrsAny = &ValidationErrors{}
)

// ValidationError is error returned for any validation errors.
// WARNING: This error will be printed to the user as is!
type ValidationError struct {
	msg string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		msg: msg,
	}
}

func NewValidationErrorf(format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		msg: fmt.Sprintf(format, args...),
	}
}

func (e *ValidationError) Error() string {
	return e.msg
}

func (e *ValidationError) Is(target error) bool {
	// If the caller is checking for any ValidationError, return true
	if errors.Is(target, ErrAny) {
		return true
	}

	// ensure it's the correct type
	err := &ValidationError{}
	if !errors.As(target, &err) {
		return false
	}

	// only the same if the message is the same
	return e.msg == err.msg
}

type ValidationErrors struct {
	errs map[string]error
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		errs: map[string]error{},
	}
}

func NewValidationErrorsKey(key string, msg string) *ValidationErrors {
	return &ValidationErrors{
		errs: map[string]error{
			key: NewValidationError(msg),
		},
	}
}

func (e *ValidationErrors) Is(target error) bool {
	// If the caller is checking for any ValidationError, return true
	if errors.Is(target, ErrsAny) {
		return true
	}

	// ensure it's the correct type
	err := &ValidationErrors{}
	if !errors.As(target, &err) {
		return false
	}

	// only the same if the message is the same
	return e.Error() == err.Error()

}

func (e *ValidationErrors) AddValidationError(key string, err error) {
	e.errs[key] = err
}

func (e *ValidationErrors) Errors() map[string]string {
	errMsg := map[string]string{}
	for k, err := range e.errs {
		errMsg[k] = err.Error()
	}
	return errMsg
}

func (e *ValidationErrors) Error() string {
	err := e.Errors()
	errMsg, _ := json.Marshal(err)
	return string(errMsg)
}

func (e *ValidationErrors) HasError() bool {
	return len(e.errs) > 0
}
