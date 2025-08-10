package domain

type ResourceAlreadyExists struct {
	Message string
}

func (err ResourceAlreadyExists) Error() string {
	return err.Message
}

func NewLoginAlreadyExists(message string) *ResourceAlreadyExists {
	return &ResourceAlreadyExists{Message: message}
}

type ResourceNotFound struct {
	Message string
}

func (err ResourceNotFound) Error() string {
	return err.Message
}

func NewResourceNotFound(message string) *ResourceNotFound {
	return &ResourceNotFound{Message: message}
}

type ProblemError struct {
	Message string
	inner   error
}

func (err *ProblemError) Error() string {
	return err.Message
}

func (err *ProblemError) Inner() error {
	return err.inner
}

func NewProblemError(message string, inner error) *ProblemError {
	return &ProblemError{Message: message, inner: inner}
}
