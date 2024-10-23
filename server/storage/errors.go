package storage

import "errors"

var (
	ErrStorageKeyNotFound    = errors.New("provided key was not found")
	ErrStorageFileWrite      = errors.New("error when writing data")
	ErrStorageFileRead       = errors.New("error when reading data")
	ErrStorageJSONMarshall   = errors.New("error when marshalling JSON")
	ErrStorageJSONUnmarshall = errors.New("error when unmarshalling JSON")
)
