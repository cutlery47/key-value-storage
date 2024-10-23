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

// creating entry from json
func FromJSON(jsonEntry []byte) (InEntry, error) {
	entry := InEntry{}
	if err := json.Unmarshal(jsonEntry, &entry); err != nil {
		log.Println(err)
		return entry, ErrJSONUnmarshall
	}
	return entry, nil
}

// entry, provided by the storage
type OutEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// converting entry to json
func ToJSON(entry OutEntry) ([]byte, error) {
	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
		return []byte{}, ErrJSONMarshall
	}

	return jsonEntry, nil
}

type Value struct {
	// essentially a value of the key
	Data string `json:"data"`
	// time info
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
