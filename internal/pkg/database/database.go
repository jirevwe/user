package database

import (
	"database/sql"

	_ "github.com/jirevwe/user/internal/pkg/database/postgres"
	"github.com/jirevwe/user/internal/pkg/database/sqlite3"
)

type Database interface {
	GetDB() *sql.DB
	FindAll(q string, out []interface{}) error
}

func New() (Database, error) {
	db, err := sqlite3.NewDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
