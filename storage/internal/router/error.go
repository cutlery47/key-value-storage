package router

import (
	"log"
	"net/http"

	"github.com/cutlery47/key-value-storage/storage/internal/storage"
)

var errStatus = map[error]int{
	storage.ErrKeyNotFound:      http.StatusNotFound,
	storage.ErrKeyAlreadyExists: http.StatusBadRequest,
}

var errMessage = map[error]string{
	storage.ErrKeyNotFound:      storage.ErrKeyNotFound.Error(),
	storage.ErrKeyAlreadyExists: storage.ErrKeyAlreadyExists.Error(),
}

type errHandler struct{}

func (h errHandler) Handle(err error) (status int, msg string) {
	log.Println("error occured:", err.Error())

	status = 500
	msg = "internal server error"

	mapStatus, ok := errStatus[err]
	if ok {
		status = mapStatus
	}

	mapMsg, ok := errMessage[err]
	if ok {
		msg = mapMsg
	}

	return status, msg
}
