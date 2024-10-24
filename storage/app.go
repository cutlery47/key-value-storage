package storage

import (
	"github.com/cutlery47/key-value-storage/storage/internal/router"
	"github.com/cutlery47/key-value-storage/storage/internal/service"
	"github.com/cutlery47/key-value-storage/storage/internal/storage"
	"github.com/cutlery47/key-value-storage/storage/server"
)

func Run() {
	ls := storage.NewLocalStorage("data")
	se := service.New(ls)
	rt := router.New(se)
	serv := server.New(rt.Handler())

	serv.Run()
}
