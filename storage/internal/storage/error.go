package storage

import "errors"

var (
	ErrKeyNotFound      = errors.New("no data was found by provided key")
	ErrKeyAlreadyExists = errors.New("provided key already exists")
	ErrFileWrite        = errors.New("error when writing data")
	ErrFileRead         = errors.New("error when reading data")
	ErrJSONMarshall     = errors.New("error when marshalling JSON")
	ErrJSONUnmarshall   = errors.New("error when unmarshalling JSON")
)
