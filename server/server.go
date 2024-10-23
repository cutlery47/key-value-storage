package server

import (
	"fmt"
	"time"

	"github.com/cutlery47/key-value-storage/server/storage"
)

func Run() {

}

func Test() {
	entry1 := storage.Entry{
		Key: "somekey1",
		Value: storage.Value{
			Data:      "someval1",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	entry2 := storage.Entry{
		Key: "somekey2",
		Value: storage.Value{
			Data:      "someval2",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	ls := storage.NewLocalStorage("data")

	if err := ls.Create(entry1); err != nil {
		fmt.Println(err.Error())
	}

	if err := ls.Create(entry2); err != nil {
		fmt.Println(err.Error())
	}

	readEntry, err := ls.Read("somekey1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(*readEntry)
	}

	readEntry2, err := ls.Read("somekey5")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(*readEntry2)
	}

	updEntry1 := storage.Entry{
		Key: "somekey1",
		Value: storage.Value{
			Data:      "someval2",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	if err := ls.Update(updEntry1); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("updated")
	}

	updEntry2 := storage.Entry{
		Key: "somekey4",
		Value: storage.Value{
			Data:      "someval2",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		},
	}

	if err := ls.Update(updEntry2); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("updated")
	}

	delKey1 := storage.Key("somekey3")
	delKey2 := storage.Key("somekey5")

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

	go ls.Cleanup(5*time.Second, finChan, sigChan, errSigChan)

	<-sigChan
}
