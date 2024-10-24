package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	Create(entry InEntry) error
	Read(key string) (*OutEntry, error)
	Update(entry InEntry) error
	Delete(key string) error
}

type LocalStorage struct {
	file fileHandler
	mu   *sync.Mutex
	log  *logrus.Logger
}

func NewLocalStorage(filepath string, log *logrus.Logger) *LocalStorage {
	file := fileHandler{
		filepath: filepath,
	}

	ls := &LocalStorage{
		file: file,
		mu:   &sync.Mutex{},
		log:  log,
	}

	// default clanup time - 10 seconds
	go ls.Cleanup(10 * time.Second)

	return ls
}

func (ls *LocalStorage) Create(entry InEntry) error {
	// reading data from a file
	data, err := ls.file.read()
	if err != nil {
		return err
	}

	// check if provided value is already stored
	if _, ok := (*data)[entry.Key]; ok {
		return ErrKeyAlreadyExists
	} else {
		(*data)[entry.Key] = entry.Value
	}

	// updating file with new data
	if err := ls.file.flush(*data); err != nil {
		return err
	}

	return nil
}

func (ls *LocalStorage) Read(key string) (*OutEntry, error) {
	// reading data from a file
	data, err := ls.file.read()
	if err != nil {
		return nil, err
	}

	// retrieve data
	val, ok := (*data)[key]
	// check if key exists
	if !ok {
		return nil, ErrKeyNotFound
	}

	return &OutEntry{Value: val.Data, Key: key}, nil
}

func (ls *LocalStorage) Update(entry InEntry) error {
	// reading data from a file
	data, err := ls.file.read()
	if err != nil {
		return err
	}

	// check if provided key is already stored
	if _, ok := (*data)[entry.Key]; !ok {
		return ErrKeyNotFound
	} else {
		(*data)[entry.Key] = entry.Value
	}

	// updating file with new data
	if err := ls.file.flush(*data); err != nil {
		return err
	}

	return nil
}

func (ls *LocalStorage) Delete(key string) error {
	// reading data from a file
	data, err := ls.file.read()
	if err != nil {
		return err
	}

	// check if data exists
	if _, ok := (*data)[key]; !ok {
		return ErrKeyNotFound
	} else {
		delete(*data, key)
	}

	// updating file with new data
	if err := ls.file.flush(*data); err != nil {
		return err
	}

	return nil
}

func (ls *LocalStorage) Cleanup(cooldown time.Duration) {
	for {
		// clean expired data each cooldown-amount seconds
		<-time.After(cooldown)

		now := time.Now()

		ls.log.WithFields(logrus.Fields{
			"status": "started",
			"at":     now,
		}).Info()

		// locking up any I/O on file until cleanup is over
		ls.mu.Lock()

		fileData, err := ls.file.read()
		if err != nil {
			// errSigChan <- err
			return
		}

		// iterating over file data and calculating
		// whether the expiration time exceedes current time
		for k, v := range *fileData {
			if v.ExpiresAt.Before(now) {
				ls.log.WithFields(logrus.Fields{
					"status": "found",
					"data":   fmt.Sprintf("key=%v: value=%v \n", k, v.Data),
					"at":     now,
				}).Info()
				// deleing entry
				delete(*fileData, k)
			}
		}

		if err := ls.file.flush(*fileData); err != nil {
			// errSigChan <- err
			return
		}

		ls.mu.Unlock()

		ls.log.WithFields(logrus.Fields{
			"status": "ended",
			"at":     now,
		}).Info()

	}
}

// abstraction over file I/O
type fileHandler struct {
	mu       *sync.Mutex
	filepath string
}

// runtime data storage
type data map[string]Value

// reads data from a file by specified filename
// converts raw JSON data into a map
func (f fileHandler) read() (*data, error) {
	byteData, err := os.ReadFile(f.filepath)
	if err != nil {
		log.Println("error when reading data:", err)
		return nil, ErrFileRead
	}

	// if our file is completely empty,
	// we add empty JSON parentheses,
	// so that unmarshaller doesn't throw UEOF errors
	if len(byteData) == 0 {
		os.WriteFile(f.filepath, []byte("{}"), 0666)
		byteData, _ = os.ReadFile(f.filepath)
	}

	currData := &data{}

	if err := json.Unmarshal(byteData, currData); err != nil {
		log.Println("error when unmarshalling data:", err)
		return nil, ErrJSONUnmarshall
	}

	return currData, nil
}

// converts incoming data into JSON
// pushes updates to a file by specified filename
func (f fileHandler) flush(newData data) error {
	newJsonData, err := json.Marshal(newData)
	if err != nil {
		log.Println("error when marshalling data:", err)
		return ErrJSONMarshall
	}

	if err := os.WriteFile(f.filepath, newJsonData, 0666); err != nil {
		log.Println("error when writing data:", err)
		return ErrFileWrite
	}

	return nil
}
