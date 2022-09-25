package database

import (
	_ "github.com/jirevwe/user/internal/pkg/database/postgres"
	"github.com/jirevwe/user/internal/pkg/database/sqlite3"
	"github.com/jirevwe/user/internal/pkg/services"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetDB() *sqlx.DB
	GetUserService() *services.UserService
}

func New() (Database, error) {
	db, err := sqlite3.NewDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
