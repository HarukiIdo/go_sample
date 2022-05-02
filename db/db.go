package db

import "database/sql"

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("splite3", path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil

}
