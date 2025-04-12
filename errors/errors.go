package errors

import (
	std_errors "errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func New(msg string) error {
	return errors.New(msg)
}

func Newf(msg string, args ...any) error {
	return fmt.Errorf(msg, args)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Wrap(err error, msg string) error {
	return Join(err, New(msg))
}

func Append(err error, errs ...error) *multierror.Error {
	return multierror.Append(err, errs...)
}

func Flatten(err error) error {
	return multierror.Flatten(err)
}

func Prefix(err error, prefix string) error {
	return multierror.Prefix(err, prefix)
}

func Join(errs ...error) error {
	return std_errors.Join(errs...)
}

func Cause(err error) error {
	cause := Unwrap(err)
	if cause == nil {
		return err
	}
	return Cause(cause)
}

// Unwrap uses causer to return the next error in the chain or nil.
// This goes one-level deeper, whereas Cause goes as far as possible
func Unwrap(err error) error {
	type causer interface {
		Cause() error
	}
	if unErr, ok := err.(causer); ok {
		return unErr.Cause()
	}
	return nil
}
