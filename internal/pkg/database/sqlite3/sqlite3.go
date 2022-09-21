package sqlite3

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const pkgName = "sqlite3"

type Sqlite struct {
	db *sql.DB
}

func NewDB() *Sqlite {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("[%s]: failed to open database - %v", pkgName, err)
	}

	return &Sqlite{db: db}
}

func (s *Sqlite) GetDB() *sql.DB {
	return s.db
}

// FindAll fetches all the entries from a query q
// and writes them into the out parameter
func (s *Sqlite) FindAll(q string, out []interface{}) error {
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
