package dto

type ValidationError struct {
	inner error
}

func NewValidationError(inner error) *ValidationError {
	return &ValidationError{inner: inner}
}

func (e *ValidationError) Error() string {
	return e.inner.Error()
}
