package Database

import "database/sql"

type Store struct{
	*Queries
	db *sql.DB
}

func ServerConn(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}