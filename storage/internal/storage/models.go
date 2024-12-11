package storage

import (
	"encoding/json"
	"log"
	"time"
)

// entry, provided to the storage
type Entry struct {
	Key   Key `json:"key"`
	Value Value
}

type Key string

type Value struct {
	// essentially a value of the key
	Data string `json:"data"`
	// time info
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// basiaclly a default InEntry constructor
func EntryFromData(key, value string, updatedAt, expiresAt time.Time) Entry {
	return Entry{
		Key: Key(key),
		Value: Value{
			Data:      value,
			UpdatedAt: updatedAt,
			ExpiresAt: expiresAt,
		},
	}
}

// converting entry to json
func (entry Entry) ToJSON() ([]byte, error) {
	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
		return []byte{}, ErrJSONMarshall
	}

	return jsonEntry, nil
}
