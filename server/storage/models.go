package storage

import "time"

type Key string

type Value struct {
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Entry struct {
	Key   Key
	Value Value
}
