package client

import "errors"

var (
	ErrParseDuration error = errors.New("couldn't parse provided duration")
	ErrEmptyKey      error = errors.New("key shouldn't be of length 0")
	ErrEmptyVal      error = errors.New("value shouldn't be of length 0")
	ErrOpUnsupported error = errors.New("operation is not supported")
	ErrOpNotProvided error = errors.New("operation is not provided")
)
