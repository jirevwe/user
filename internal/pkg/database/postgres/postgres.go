package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const pkgName = "postgres"

type Postgres struct {
	db *sql.DB
}

func NewDB() (*Postgres, error) {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/user")
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to open database - %v", pkgName, err)
	}

	return &Postgres{db: db}, nil
}

func (s *Postgres) GetDB() *sql.DB {
	return s.db
}

// FindAll fetches all the entries from a query q
// and writes them into the out parameter
func (s *Postgres) FindAll(q string, out []interface{}) error {
	rows, err := s.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var o interface{}
		err := rows.Scan(&o)
		if err != nil {
			return err
		}

		out = append(out, o)
	}

	// check if the cursor had any error
	return rows.Err()
}
