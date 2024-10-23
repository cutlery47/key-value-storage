package service

import (
	"time"

	"github.com/cutlery47/key-value-storage/server/storage"
)

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) Add(key, value, expires_at string) error {
	timeUpdatedAt := time.Now()
	timeExpiresAt, err := time.Parse("", expires_at)
	if err != nil {
		return err
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

func (s Service) Set(key, value, expires_at string) error {
	timeUpdateddAt := time.Now()
	timeExpiresAt, err := time.Parse("", expires_at)
	if err != nil {
		return err
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

func (s Service) Get(key string) (string, error) {
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

func (s Service) Delete(key string) error {
	return s.storage.Delete(key)
}
