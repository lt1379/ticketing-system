package repository_test

import (
	"github.com/lts1379/ticketing-system/infrastructure/persistence"
	"testing"
)

func Test_GetById(t *testing.T) {
	db, _ := persistence.NewNativeDb()
	db.Close()
}
