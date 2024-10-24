package service

import (
	"time"

	"github.com/cutlery47/key-value-storage/storage/internal/storage"
)

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

	entry := storage.InEntry{
		Key: key,
		Value: storage.Value{
			Data:      value,
			UpdatedAt: timeUpdatedAt,
			ExpiresAt: timeExpiresAt,
		},
	}

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

	entry := storage.InEntry{
		Key: key,
		Value: storage.Value{
			Data:      value,
			UpdatedAt: timeUpdateddAt,
			ExpiresAt: timeExpiresAt,
		},
	}

	return s.storage.Update(entry)
}

func (s *Service) Get(key string) (string, error) {
	entry, err := s.storage.Read(key)
	if err != nil {
		return "", err
	}

	jsonEntry, err := storage.ToJSON(*entry)
	if err != nil {
		return "", err
	}

	return string(jsonEntry), nil
}

func (s *Service) Delete(key string) error {
	return s.storage.Delete(key)
}
