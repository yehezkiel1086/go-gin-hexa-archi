package domain

import "errors"

var (
	ErrInternal = errors.New("internal server error")
	ErrNotFound = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidQuery = errors.New("invalid queries")
	ErrConflictingData = errors.New("conflicting data")
	ErrInvalidIDParam = errors.New("invalid id parameter")
	ErrUserNotFound = errors.New("user is not found")
)
