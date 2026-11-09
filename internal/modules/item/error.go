package item

import "errors"

var (
	ErrInvalidTags      = errors.New("invalid tags")
	ErrTagNotFound      = errors.New("tag not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotFound         = errors.New("not found")
	ErrImageNotUploaded = errors.New("image not uploaded")
)
