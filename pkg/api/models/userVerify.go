package models

import "database/sql"

// UserVerify models user verification data
type UserVerify struct {
	ID       int
	UserID   int
	NumberID int
	Code     string
}

// NewUserVerify creates a new userVerify entry in DB
func (uv *UserVerify) NewUserVerify(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO verifyUser(userID, numberID, code) VALUES(?, ?, ?)`,
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

	return id, nil
}

// GetUserVerifyByID returns user verification data using
// given ID
func (uv *UserVerify) GetUserVerifyByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT userID, numberID, code FROM userVerify WHERE id = ?`,
		uv.ID,
	).Scan(&uv.UserID, &uv.NumberID, &uv.Code)
}

// GetUserVerifyByUserID returns user verification data using
// given UserID
func (uv *UserVerify) GetUserVerifyByUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, numberID, code FROM userVerify WHERE userID = ?`,
		uv.UserID,
	).Scan(&uv.ID, &uv.NumberID, &uv.Code)
}

// DeleteUserVerifyByID deletes user verification entry by using
// given ID
func (uv *UserVerify) DeleteUserVerifyByID(tx *sql.Tx) (int64, error) {
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

// DeleteUserVerifyByUserID deletes user verification entry by using
// given userID
func (uv *UserVerify) DeleteUserVerifyByUserID(tx *sql.Tx) (int64, error) {
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
