package filecsv

import (
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"os"
)

func NewFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while open file")
		return nil, err
	}

	return file, nil
}
