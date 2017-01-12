package implement

type Error struct {
	Code int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

const (
	InvalidLengthErrCode = iota
)

var ErrInvalidLengths = NewError(InvalidLengthErrCode, "Interfaces to read in must be the same length as the number of implementor names.")
