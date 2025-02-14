package errorsPackage

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found") // only storage and service
	ErrUserExists         = errors.New("user already exists")
	ErrAppNotFound        = errors.New("app not found")
	ErrRoleNotFound       = errors.New("role not found")
	ErrInvalidCredentials = errors.New("login or password is incorrect")
)
