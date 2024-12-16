package service

import (
	"time"

	"github.com/cutlery47/key-value-storage/storage/internal/storage"
)

// handles and transforms incoming request data
// passes entries down to the storage layer
type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Add(key, value, expiresAt string) error {
	var timeExpiresAt time.Time
	timeUpdatedAt := time.Now()

	// if ttl was not provided - set default to 24 hours
	if len(expiresAt) == 0 {
		timeExpiresAt = time.Now().Add(time.Hour * 24)
	} else {
		parsed, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return err
		}
		timeExpiresAt = parsed
	}

	entry := storage.EntryFromData(key, value, timeUpdatedAt, timeExpiresAt)

	return s.storage.Create(entry)
}

func (s *Service) Set(key, value, expiresAt string) error {
	var timeExpiresAt time.Time
	timeUpdateddAt := time.Now()

	if len(expiresAt) != 0 {
		parsed, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return err
		}
		timeExpiresAt = parsed
	}

	entry := storage.EntryFromData(key, value, timeUpdateddAt, timeExpiresAt)

	return s.storage.Update(entry)
}

func (s *Service) Get(key string) (string, error) {
	entry, err := s.storage.Read(storage.Key(key))
	if err != nil {
		return "", err
	}

	jsonEntry, err := entry.ToJSON()
	if err != nil {
		return "", err
	}

	return string(jsonEntry), nil
}

func (s *Service) Delete(key string) error {
	return s.storage.Delete(storage.Key(key))
}
