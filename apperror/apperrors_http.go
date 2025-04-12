package apperror

import (
	"errors"
)

func FromError(err error) (*AppError, bool) {
	type httperror interface{ HTTPError() *AppError }

	if gs, ok := err.(httperror); ok {
		httpError := gs.HTTPError()
		if httpError == nil {
			return nil, false
		}

		return httpError, true
	}

	var gs httperror
	if errors.As(err, &gs) {
		httpError := gs.HTTPError()
		if httpError == nil {
			return nil, false
		}

		return httpError, true
	}

	return nil, false
}

func (e *AppError) HTTPError() *AppError {
	return newAppError(
		e.Type,
		e.SystemCode,
		WithCode(e.Code),
		WithMessage(e.Message),
		WithDomain(e.Domain),
		WithErr(e.Err),
		WithFields(e.Fields),
		WithTraceID(e.TraceID),
	)
}
