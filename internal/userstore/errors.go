package userstore

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrDuplicateUser = errors.New("User already exists")
)