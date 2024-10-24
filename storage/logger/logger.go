package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func NewJsonFile(filepath string, level logrus.Level) (*logrus.Logger, error) {
	logger := logrus.New()

	fd, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	logger.SetOutput(fd)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)

	return logger, nil
}
