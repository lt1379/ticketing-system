package persistence

import (
	"database/sql"
	"fmt"
	"my-project/infrastructure/configuration"
	"time"
)

func NewNativeDb() (*sql.DB, error) {
	cfg := configuration.C.Database.Db

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}
