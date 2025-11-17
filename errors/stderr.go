package errors

import "errors"

func New(text string) error {
	return errors.New(text)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
