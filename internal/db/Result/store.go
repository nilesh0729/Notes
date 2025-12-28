package Database

import "database/sql"

type Store interface{
	Querier
}

type RealStore struct{
	*Queries
	db *sql.DB
}

func ServerConn(db *sql.DB) Store {
	return &RealStore{
		db:      db,
		Queries: New(db),
	}
}