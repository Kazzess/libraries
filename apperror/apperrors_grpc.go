package apperror

import (
	"strconv"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (e *AppError) GRPCStatus() *status.Status {
	var code codes.Code

	switch e.Type {
	case errInternalSystemCode:
		code = codes.Internal
	case errBadRequestCode, errValidationCode:
		code = codes.InvalidArgument
	case errNotFoundCode:
		code = codes.NotFound
	case errUnauthorizedCode:
		code = codes.Unauthenticated
	case errForbiddenCode:
		code = codes.PermissionDenied
	case errConditionFailedCode:
		code = codes.FailedPrecondition
	default:
		code = codes.Unknown
	}

	st := status.New(code, e.Message)

	var details []proto.Message

	errorDetails := e.errorDetails()
	details = append(details, errorDetails)

	if e.Fields != nil {
		fieldDetails := e.fieldDetails()
		details = append(details, fieldDetails)
	}

	withDetails, err := st.WithDetails(details...)
	if err != nil {
		return st
	}

	return withDetails
}

func (e *AppError) fieldDetails() proto.Message {
	br := &errdetails.BadRequest{}
	for k, v := range e.Fields {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       k,
			Description: v,
		})
	}

	return br
}

func (e *AppError) errorDetails() proto.Message {
	br := &errdetails.ErrorInfo{
		Reason:   e.SystemCode + "-" + strconv.FormatInt(int64(e.Code), 10),
		Domain:   e.Domain,
		Metadata: nil,
	}

	return br
}
