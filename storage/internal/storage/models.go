package storage

import (
	"encoding/json"
	"log"
	"time"
)

// entry, provided to the storage
type InEntry struct {
	Key   string `json:"key"`
	Value Value
}

type Value struct {
	// essentially a value of the key
	Data string `json:"data"`
	// time info
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// basiaclly a default InEntry constructor
func EntryFromData(key, value string, updatedAt, expiresAt time.Time) InEntry {
	return InEntry{
		Key: key,
		Value: Value{
			Data:      value,
			UpdatedAt: updatedAt,
			ExpiresAt: expiresAt,
		},
	}
}

// entry, provided by the storage
type OutEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// converting entry to json
func (entry OutEntry) ToJSON() ([]byte, error) {
	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
		return []byte{}, ErrJSONMarshall
	}

	return jsonEntry, nil
}
