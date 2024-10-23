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

func (s *Service) Add(key, value, expiresAt string) error {
	timeUpdatedAt := time.Now()
	timeExpiresAt, err := s.getExpirationTime(expiresAt)
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

func (s *Service) Set(key, value, expiresAt string) error {
	timeUpdateddAt := time.Now()
	timeExpiresAt, err := s.getExpirationTime(expiresAt)
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

func (s *Service) getExpirationTime(expiresAt string) (time.Time, error) {
	var timeExpiresAt time.Time

	// if expiration time was not provided - set default
	if len(expiresAt) == 0 {
		timeExpiresAt = time.Now().Add(time.Hour * 24)
	} else {
		parsedTime, err := time.Parse(time.RFC1123Z, expiresAt)
		if err != nil {
			return time.Time{}, err
		}
		timeExpiresAt = parsedTime
	}

	return timeExpiresAt, nil
}
