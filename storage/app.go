package storage

import (
	"log"

	"github.com/cutlery47/key-value-storage/storage/internal/router"
	"github.com/cutlery47/key-value-storage/storage/internal/service"
	"github.com/cutlery47/key-value-storage/storage/internal/storage"
	"github.com/cutlery47/key-value-storage/storage/logger"
	"github.com/cutlery47/key-value-storage/storage/server"
	"github.com/sirupsen/logrus"
)

func Run() {
	// request logger
	reqLog, err := logger.NewJsonFile("logger/logs/requests.log", logrus.InfoLevel)
	if err != nil {
		log.Fatal("couldn't configure request logger:", err)
	}

	// cleanup logger
	cleLog, err := logger.NewJsonFile("logger/logs/cleanup.log", logrus.InfoLevel)
	if err != nil {
		log.Fatal("couldn't configure cleanup logger", err)
	}

	// error logger
	errLog, err := logger.NewJsonFile("logger/logs/error.log", logrus.ErrorLevel)
	if err != nil {
		log.Fatal("couldn't configure error logger", err)
	}

	ls := storage.NewLocalStorage("storage.data", cleLog)
	se := service.New(ls)
	rt := router.New(se, reqLog, errLog)
	serv := server.New(rt.Handler())

	serv.Run()
}
