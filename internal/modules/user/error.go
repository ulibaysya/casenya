package user

import "errors"

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrNotFound            = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrAlreadyExists       = errors.New("already exists")
	ErrSessionNotFound     = errors.New("session not found")
	ErrItemNotActive       = errors.New("item not active")
	ErrItemNotFound        = errors.New("item not found")
)
