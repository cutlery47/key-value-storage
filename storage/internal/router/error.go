package router

import (
	"net/http"
	"time"

	"github.com/cutlery47/key-value-storage/storage/internal/storage"
	"github.com/sirupsen/logrus"
)

var errStatus = map[error]int{
	storage.ErrKeyNotFound:      http.StatusNotFound,
	storage.ErrKeyAlreadyExists: http.StatusBadRequest,
}

type errHandler struct {
	errLog *logrus.Logger
}

func (h errHandler) Handle(err error) (status int, msg string) {
	status = 500
	msg = "internal server error"

	// if error is not internal - map it to specific status
	// else return 500 and log out the error
	mapStatus, ok := errStatus[err]
	if ok {
		status = mapStatus
		msg = err.Error()
	} else {
		h.errLog.WithFields(
			logrus.Fields{
				"time":  time.Now(),
				"error": err.Error(),
			},
		).Error()
	}

	return status, msg
}
