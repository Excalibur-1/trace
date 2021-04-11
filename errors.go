package trace

import (
	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedFormat occurs when the `format` passed to Tracer.Inject() or
	// Tracer.Extract() is not recognized by the Tracer implementation.
	ErrUnsupportedFormat = errors.New("trace: Unknown or unsupported Inject/Extract format")

	// ErrTraceNotFound occurs when the `carrier` passed to
	// Tracer.Extract() is valid and uncorrupted but has insufficient
	// information to extract a Trace.
	ErrTraceNotFound = errors.New("trace: Trace not found in Extract carrier")

	// ErrInvalidTrace errors occur when Tracer.Inject() is asked to
	// operate on a Trace which it is not prepared to handle (for
	// example, since it was created by a different tracer implementation).
	ErrInvalidTrace = errors.New("trace: Trace type incompatible with tracer")

	// ErrInvalidCarrier errors occur when Tracer.Inject() or Tracer.Extract()
	// implementations expect a different type of `carrier` than they are
	// given.
	ErrInvalidCarrier = errors.New("trace: Invalid Inject/Extract carrier")

	// ErrTraceCorrupted occurs when the `carrier` passed to
	// Tracer.Extract() is of the expected type but is corrupted.
	ErrTraceCorrupted = errors.New("trace: Trace data corrupted in Extract carrier")

	errEmptyTracerString = errors.New("trace: cannot convert empty string to span context")

	errInvalidTracerString = errors.New("trace: string does not match span context string format")
)
