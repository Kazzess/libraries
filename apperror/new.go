package apperror

func NewInternalError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errInternalSystemCode,
		systemCode,
		options...,
	)
}

func NewBadRequestError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errBadRequestCode,
		systemCode,
		options...,
	)
}

func NewValidationError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errValidationCode,
		systemCode,
		options...,
	)
}

func NewNotFoundError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errNotFoundCode,
		systemCode,
		options...,
	)
}

func NewUnauthorizedError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errUnauthorizedCode,
		systemCode,
		options...,
	)
}

func NewForbiddenError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errForbiddenCode,
		systemCode,
		options...,
	)
}

func NewConditionFailedError(systemCode string, options ...Option) *AppError {
	return newAppError(
		errConditionFailedCode,
		systemCode,
		options...,
	)
}

type Option func(*AppError)

// WithErr Option setter for Err.
func WithErr(err error) Option {
	return func(ae *AppError) {
		ae.Err = err
	}
}

// WithMessage Option setter for Message.
func WithMessage(message string) Option {
	return func(ae *AppError) {
		ae.Message = message
	}
}

// WithDomain Option setter for Domain.
func WithDomain(domain string) Option {
	return func(ae *AppError) {
		ae.Domain = domain
	}
}

// WithCode Option setter for Code.
func WithCode(code uint32) Option {
	return func(ae *AppError) {
		ae.Code = code
	}
}

// WithFields Option setter for Fields.
func WithFields(fields ErrorFields) Option { // Assuming 'ErrorFields' is a known type in your code
	return func(ae *AppError) {
		ae.Fields = fields
	}
}

// WithTraceID Option setter for TraceID.
func WithTraceID(traceID string) Option {
	return func(ae *AppError) {
		ae.TraceID = traceID
	}
}
