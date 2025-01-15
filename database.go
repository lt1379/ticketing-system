package main

import (
	"database/sql"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"github.com/lts1379/ticketing-system/infrastructure/persistence"
)

func InitiateDatabase() (*sql.DB, *sql.DB, error) {
	var err error

	db, err := persistence.NewNativeDb()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Cannot connect to the local database")
		return nil, nil, err
	}

	postgres, err := persistence.NewPostgreSQLDb()
	if err != nil {
		return nil, nil, err
	}

	return db, postgres, err
}
