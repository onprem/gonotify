package models

import (
	"database/sql"
)

// NumberVerify models data required to verify a phone number in DB
type NumberVerify struct {
	ID       int
	NumberID int
	Code     string
}

// New adds a new number verification data in database
func (nv *NumberVerify) New(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO numberVerify(numberID, code) VALUES (?, ?)`,
		nv.NumberID,
		nv.Code,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	nv.ID = int(id)
	return id, nil
}

// GetByID returns a user given it's ID
func (nv *NumberVerify) GetByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT numberID, code FROM numberVerify WHERE id = ?`,
		nv.ID,
	).Scan(&nv.NumberID, &nv.Code)
}

// GetByNumberID returns a user given it's ID
func (nv *NumberVerify) GetByNumberID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, code FROM numberVerify WHERE numberID = ?`,
		nv.NumberID,
	).Scan(&nv.ID, &nv.Code)
}

// DeleteByID deletes the number verification data using given ID
func (nv *NumberVerify) DeleteByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(`DELETE FROM numberVerify WHERE id = ?`, nv.ID)
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}
