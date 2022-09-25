package postgres

import (
	"fmt"

	"github.com/jirevwe/user/internal/pkg/services"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const pkgName = "postgres"

type Postgres struct {
	dbx *sqlx.DB
}

func NewDB() (*Postgres, error) {
	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost/user")
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to open database - %v", pkgName, err)
	}

	return &Postgres{dbx: db}, nil
}

func (p *Postgres) GetDB() *sqlx.DB {
	return p.dbx
}

func (p *Postgres) GetUserService() *services.UserService {
	return &services.UserService{DB: p.dbx}
}
