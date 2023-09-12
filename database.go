package main

import (
	"database/sql"
	"my-project/infrastructure/logger"
	"my-project/infrastructure/persistence"
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
