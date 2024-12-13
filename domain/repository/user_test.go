package repository_test

import (
	"my-project/infrastructure/persistence"
	"testing"
)

func Test_GetById(t *testing.T) {
	db, _ := persistence.NewNativeDb()
	db.Close()
}
