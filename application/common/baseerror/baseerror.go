package baseerror
type (
	BaseError struct {
		message string
	}
)

func NewBaseError(message string) *BaseError {
	return &BaseError{message: message}
}

func (e *BaseError) Error() string {

	return e.message
}