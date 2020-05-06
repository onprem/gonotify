package models

import (
	"database/sql"
)

// User models a user in database
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Hash     string `json:"-"`
	Verified bool   `json:"verified"`
}

// InsertUser adds a new user in database
func (u *User) InsertUser(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO users(name, phone, hash, verified) VALUES (?, ?, ?, ?)`,
		u.Name,
		u.Phone,
		u.Hash,
		u.Verified,
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

// GetUserByID returns a user given it's ID
func (u *User) GetUserByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT name, phone, hash, verified FROM users WHERE id = ?`,
		u.ID,
	).Scan(&u.Name, &u.Phone, &u.Hash, &u.Verified)
}

// GetUserByPhone returns a user given it's phone number
func (u *User) GetUserByPhone(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, name, hash, verified FROM users WHERE phone = ?`,
		u.Phone,
	).Scan(&u.ID, &u.Name, &u.Hash, &u.Verified)
}

// UpdateUserByID updates a user given it's ID
func (u *User) UpdateUserByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`UPDATE users SET name = ?, phone = ?, hash = ?, verified = ? WHERE id = ?`,
		u.Name,
		u.Phone,
		u.Hash,
		u.Verified,
		u.ID,
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
