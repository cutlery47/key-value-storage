package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ImprovedStorage struct {
	cc *cache

	errLog *logrus.Logger
}

func NewImprovedStorage(filepath string, errLog *logrus.Logger) (*ImprovedStorage, error) {
	st := &ImprovedStorage{
		cc: &cache{
			sync.RWMutex{},
			make(store),
		},

		errLog: errLog,
	}

	fd, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %v", err)
	}

	if err := st.restore(fd); err != nil {
		if !errors.Is(err, ErrNothingToRestore) {
			log.Println("failed to restore state: ", err)
		}
		return nil, err
	}

	go st.flush(fd, time.Minute)

	return st, nil
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

func (st *ImprovedStorage) flush(fd *os.File, to time.Duration) {
	for {
		time.Sleep(to)

		data, err := json.Marshal(st.cc.data)
		if err != nil {
			log.Println("json.Marshal:", err)
			return
		}

		if _, err := fd.Write(data); err != nil {
			log.Println("fd.Write:", err)
		}
	}
}

// restores storage state from disk
func (st *ImprovedStorage) restore(fd *os.File) error {
	stat, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("fd.Stat: %v", err)
	}

	buf := make([]byte, stat.Size())
	if _, err := fd.Read(buf); err != nil {
		return fmt.Errorf("fd.Read: %v", err)
	}

	if err := json.Unmarshal(buf, &st.cc.data); err != nil {
		return fmt.Errorf("json.Unmarshall: %v", err)
	}

	return nil
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
