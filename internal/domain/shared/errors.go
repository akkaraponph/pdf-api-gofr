package shared

import "fmt"

// ErrInvalidInput indicates a client-side validation error.
type ErrInvalidInput struct {
	Message string
}

func (e ErrInvalidInput) Error() string  { return e.Message }
func (e ErrInvalidInput) StatusCode() int { return 400 }

// ErrProcessingFailed indicates the PDF engine failed to process a valid request.
type ErrProcessingFailed struct {
	Operation string
	Cause     error
}

func (e ErrProcessingFailed) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Cause)
}
func (e ErrProcessingFailed) StatusCode() int { return 422 }
func (e ErrProcessingFailed) Unwrap() error   { return e.Cause }

// ErrToolNotFound indicates a required external tool is not installed.
type ErrToolNotFound struct {
	Tool string
}

func (e ErrToolNotFound) Error() string {
	return fmt.Sprintf("required tool not found: %s (install poppler-utils, mupdf-tools, or ghostscript)", e.Tool)
}
func (e ErrToolNotFound) StatusCode() int { return 503 }
