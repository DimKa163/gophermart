package mediatr

type MediatrRegisterHandlerError struct {
	Message string
}

func NewMediatrRegisterHandlerError(message string) *MediatrRegisterHandlerError {
	return &MediatrRegisterHandlerError{
		Message: message,
	}
}

func (e *MediatrRegisterHandlerError) Error() string {
	return e.Message
}

type MediatrNotFoundHandlerError struct {
	Message string
}

func NewMediatrNotFoundHandlerError(message string) *MediatrNotFoundHandlerError {
	return &MediatrNotFoundHandlerError{
		Message: message,
	}
}

func (e *MediatrNotFoundHandlerError) Error() string {
	return e.Message
}

type MediatrConvertingHandlerError struct {
	Message string
}

func NewMediatrConvertingHandlerError(message string) *MediatrConvertingHandlerError {
	return &MediatrConvertingHandlerError{
		Message: message,
	}
}

func (e *MediatrConvertingHandlerError) Error() string {
	return e.Message
}
