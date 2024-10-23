package server

import (
	"fmt"
	"time"

	"github.com/cutlery47/key-value-storage/server/router"
	"github.com/cutlery47/key-value-storage/server/server"
	"github.com/cutlery47/key-value-storage/server/service"
	"github.com/cutlery47/key-value-storage/server/storage"
)

func Run() {
	storage := storage.NewLocalStorage("data")
	service := service.New(storage)
	router := router.New(service)
	serv := server.New(router.Handler())

	serv.Serve()
}

func Test() {
	entry1 := storage.InEntry{
		Key: "somekey1",
		Value: storage.Value{
			Data:      "someval1",
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	entry2 := storage.InEntry{
		Key: "somekey3",
		Value: storage.Value{
			Data:      "someval3",
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Minute),
		},
	}

	ls := storage.NewLocalStorage("data")

	if err := ls.Create(entry1); err != nil {
		fmt.Println("create 1:", err.Error())
	}

	if err := ls.Create(entry2); err != nil {
		fmt.Println("create 2:", err.Error())
	}

	readEntry, err := ls.Read("somekey1")
	if err != nil {
		fmt.Println("read1:", err.Error())
	} else {
		fmt.Println(*readEntry)
	}

	readEntry2, err := ls.Read("somekey5")
	if err != nil {
		fmt.Println("read2:", err.Error())
	} else {
		fmt.Println(*readEntry2)
	}

	updEntry1 := storage.InEntry{
		Key: "somekey1",
		Value: storage.Value{
			Data:      "someval2",
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Minute),
		},
	}

	if err := ls.Update(updEntry1); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("updated")
	}

	updEntry2 := storage.InEntry{
		Key: "somekey4",
		Value: storage.Value{
			Data:      "someval2",
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	if err := ls.Update(updEntry2); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("updated")
	}

	delKey1 := "somekey3"
	delKey2 := "somekey5"

	if err := ls.Delete(delKey1); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("deleted")
	}

	if err := ls.Delete(delKey2); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("deleted")
	}

	finChan := make(chan byte)
	sigChan := make(chan byte)
	errSigChan := make(chan error)

	go ls.Cleanup(time.Second, finChan, sigChan, errSigChan)

	select {
	case <-sigChan:
	case err := <-errSigChan:
		fmt.Println(err)
	}
}
