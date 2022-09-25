package migrator

import (
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dbx *sqlx.DB
	src migrate.MigrationSource
}

func New(d database.Database) *Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "sql",
	}
	return &Migrator{dbx: d.GetDB(), src: migrations}
}

func (m *Migrator) Up() error {
	_, err := migrate.Exec(m.dbx.DB, "sqlite3", m.src, migrate.Up)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Down() error {
	_, err := migrate.Exec(m.dbx.DB, "sqlite3", m.src, migrate.Down)
	if err != nil {
		return err
	}
	return nil
}
