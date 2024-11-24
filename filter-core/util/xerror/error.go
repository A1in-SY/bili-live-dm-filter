package xerror

type Error struct {
	code    int32
	message string
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() int32 {
	return e.code
}

func New(code int32, msg string) *Error {
	return &Error{
		code:    code,
		message: msg,
	}
}

var (
	DefaultError = New(500, "default error")
)
