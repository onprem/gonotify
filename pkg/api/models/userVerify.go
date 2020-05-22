package models

import "database/sql"

// UserVerify models user verification data
type UserVerify struct {
	ID       int
	UserID   int
	NumberID int
	Code     string
}

// New creates a new userVerify entry in DB
func (uv *UserVerify) New(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO userVerify(userID, numberID, code) VALUES(?, ?, ?)`,
		uv.UserID,
		uv.NumberID,
		uv.Code,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	uv.ID = int(id)
	return id, nil
}

// GetByID returns user verification data using
// given ID
func (uv *UserVerify) GetByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT userID, numberID, code FROM userVerify WHERE id = ?`,
		uv.ID,
	).Scan(&uv.UserID, &uv.NumberID, &uv.Code)
}

// GetByUserID returns user verification data using
// given UserID
func (uv *UserVerify) GetByUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, numberID, code FROM userVerify WHERE userID = ?`,
		uv.UserID,
	).Scan(&uv.ID, &uv.NumberID, &uv.Code)
}

// DeleteByID deletes user verification entry by using
// given ID
func (uv *UserVerify) DeleteByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`DELETE FROM userVerify WHERE id = ?`,
		uv.ID,
	)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// DeleteByUserID deletes user verification entry by using
// given userID
func (uv *UserVerify) DeleteByUserID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`DELETE FROM userVerify WHERE userID = ?`,
		uv.UserID,
	)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
