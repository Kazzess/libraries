package apperror

import (
	"context"
	"encoding/json"
	"errors"

	"go.opentelemetry.io/otel/trace"
)

type ErrorFields map[string]string

type AppError struct {
	Err        error       `json:"-"`
	Message    string      `json:"message"`
	Domain     string      `json:"domain,omitempty"`
	SystemCode string      `json:"system_code,omitempty"`
	Type       Type        `json:"type"`
	Code       uint32      `json:"code"`
	Fields     ErrorFields `json:"fields,omitempty"`
	TraceID    string      `json:"trace_id,omitempty"`
}

func (e *AppError) WithFields(fields ErrorFields) *AppError {
	e.Fields = fields

	return e
}

func newAppError(errType Type, systemCode string, options ...Option) *AppError {
	ae := &AppError{
		Type:       errType,
		SystemCode: systemCode,
	}

	for _, opt := range options {
		opt(ae)
	}

	if ae.Err == nil {
		ae.Err = errors.New(ae.Message)
	}

	if ae.Message == "" && ae.Err != nil {
		ae.Message = ae.Err.Error()
	}

	return ae
}

func (e *AppError) Error() string {
	err := e.Err.Error()

	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			err += ", " + k + " " + v
		}
	}

	return err
}

func (e *AppError) WithTrace(ctx context.Context) *AppError {
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		e.TraceID = span.TraceID().String()
	}

	return e
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}

	return bytes
}
