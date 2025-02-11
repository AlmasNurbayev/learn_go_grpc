package storage

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrAppNotFound  = errors.New("app not found")
	ErrRoleNotFound = errors.New("role not found")
)
