package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ImprovedStorage struct {
	cc *cache
	fl *file

	infoLog *logrus.Logger
	errLog  *logrus.Logger

	errChan chan<- error
}

func NewImprovedStorage(errLog, infoLog *logrus.Logger) *ImprovedStorage {
	return &ImprovedStorage{
		infoLog: infoLog,
		errLog:  errLog,
	}
}

func (st *ImprovedStorage) Create(entry Entry) error {
	if _, ok := st.cc.get(entry.Key); ok {
		return ErrKeyAlreadyExists
	}

	st.cc.put(entry)
	return nil
}

func (st *ImprovedStorage) Read(key Key) (Entry, error) {
	val, ok := st.cc.get(key)
	if !ok {
		return Entry{}, ErrKeyNotFound
	}

	return val, nil
}

func (st *ImprovedStorage) Update(entry Entry) error {
	if _, ok := st.cc.get(entry.Key); !ok {
		return ErrKeyNotFound
	}

	st.cc.put(entry)
	return nil
}

func (st *ImprovedStorage) Delete(key Key) error {
	if _, ok := st.cc.get(key); !ok {
		return ErrKeyNotFound
	}

	st.cc.del(key)
	return nil
}

func (st *ImprovedStorage) flush(to time.Duration) {
	for {
		st.cc.RLock()
		st.fl.flush(st.cc.data, st.errChan)
		st.cc.RUnlock()

		time.Sleep(to)
	}
}

// restores storage state from disk
func (st *ImprovedStorage) restore() error {
	data, err := st.fl.read()
	if err != nil {
		return err
	}

	///...
}

// in-mem storage
type store map[Key]Value

// in-mem storage
type cache struct {
	sync.RWMutex
	data store
}

func (cc *cache) put(entry Entry) {
	cc.Lock()
	cc.data[entry.Key] = entry.Value
	cc.Unlock()
}

func (cc *cache) get(key Key) (Entry, bool) {
	cc.RLock()
	val, ok := cc.data[key]
	cc.RUnlock()

	return Entry{
		Key:   key,
		Value: val,
	}, ok
}

func (cc *cache) del(key Key) {
	cc.Lock()
	delete(cc.data, key)
	cc.Unlock()
}

// abstraction over file i/o
type file struct {
	fd *os.File
}

func newFile(filepath string) (*file, error) {
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &file{
		fd: fd,
	}, nil
}

// flush data onto disk
func (f *file) flush(data store, errChan chan<- error) {
	raw, err := json.Marshal(data)
	if err != nil {
		errChan <- fmt.Errorf("json.Marshall: %v", err)
		return
	}

	if _, err = f.fd.Write(raw); err != nil {
		errChan <- fmt.Errorf("f.fd.Write: %v", err)
		return
	}
}

// read data from disk
func (f *file) read() ([]byte, error) {
	buf := []byte{}
	if _, err := f.fd.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
