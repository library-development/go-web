package web

type Error struct {
	Err string `json:"error",omitempty,omitemptykey:""`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(msg string) Error {
	return Error{Err: msg}
}

var NoError = NewError("")

var ErrInvalidSession = NewError("invalid session")
var ErrInvalidPassword = NewError("invalid password")
var ErrUserNotFound = NewError("user not found")
var ErrInvalidRegistrationCode = NewError("invalid registration code")
var ErrMethodNotSupported = NewError("method not supported")
