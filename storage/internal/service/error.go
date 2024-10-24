package service

import "errors"

var (
	ErrEmptyKey error = errors.New("empty keys are not acceptable")
)
