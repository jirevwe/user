package database

import (
	"database/sql"
)

type Database interface {
	GetDB() *sql.DB
	FindAll(q string, out []interface{}) error
}
