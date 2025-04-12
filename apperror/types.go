package apperror

type Type uint16

func (t Type) Uint32() uint32 {
	return uint32(t)
}

const (
	errInternalSystemCode Type = iota
	errBadRequestCode
	errValidationCode
	errNotFoundCode
	errUnauthorizedCode
	errForbiddenCode
	errConditionFailedCode
)
