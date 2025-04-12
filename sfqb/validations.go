package sfqb

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

var (
	ErrMinFailed     = errors.New("bad minimum value")
	ErrMaxFailed     = errors.New("bad maximum value")
	ErrMinMaxFailed  = errors.New("bad range min-max value")
	ErrEmpty         = errors.New("this params is required")
	ErrBadDateFormat = errors.New("bad date format")
)

// Min validation if value greater or equal then min
func Min(min int) QueryValidationFunc {
	return func(value interface{}) error {
		if limit, ok := value.(int); ok {
			if limit >= min {
				return nil
			}
		}
		return errors.Wrapf(ErrMinFailed, "%v", value)
	}
}

// Max validation if value lower or equal then max
func Max(max int) QueryValidationFunc {
	return func(value interface{}) error {
		if limit, ok := value.(int); ok {
			if limit <= max {
				return nil
			}
		}
		return errors.Wrapf(ErrMaxFailed, "%v", value)
	}
}

// MinMax validation if value between or equal min and max
func MinMax(min, max int) QueryValidationFunc {
	return func(value interface{}) error {
		if limit, ok := value.(int); ok {
			if min <= limit && limit <= max {
				return nil
			}
		}
		return errors.Wrapf(ErrMinMaxFailed, "%v", value)
	}
}

// NotEmpty validation if string value length more then 0
func NotEmpty() QueryValidationFunc {
	return func(value interface{}) error {
		if s, ok := value.(string); ok {
			if len(s) > 0 {
				return nil
			}
		}
		return errors.Wrapf(ErrEmpty, "%v", value)
	}
}

// DateFormat validation if string is in date format
func DateFormat(format string) QueryValidationFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return ErrBadDateFormat
		}
		return validation.Validate(s, validation.Date(format))
	}
}
