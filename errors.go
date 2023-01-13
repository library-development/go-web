package web

import "errors"

var ErrInvalidSession = errors.New("invalid session")
var ErrInvalidPassword = errors.New("invalid password")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidRegistrationCode = errors.New("invalid registration code")
