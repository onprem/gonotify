package models

import "database/sql"

// WithDB is a wrapper to run insertions and updates with DB
// It returns ID in case of insertion and rows affected in case
// of an update
func WithDB(db *sql.DB, fn func(*sql.Tx) (int64, error)) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := fn(tx)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return id, nil
}
