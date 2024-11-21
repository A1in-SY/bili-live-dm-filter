package xerror

type Error struct {
	ECode   int32
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Code() int32 {
	return e.ECode
}

func New(code int32, msg string) *Error {
	return &Error{
		ECode:   code,
		Message: msg,
	}
}
