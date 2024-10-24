package router

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func WithLogging(h http.Handler, log *logrus.Logger) http.Handler {
	logFunc := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"URI":      uri,
			"Method":   method,
			"Duration": duration,
		}).Info()
	}
	return http.HandlerFunc(logFunc)
}
