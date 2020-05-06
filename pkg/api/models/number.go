package models

import "database/sql"

type Number struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userID"`
	Phone    string `json:"phone"`
	Verified bool   `json:"verified"`
	Groups   int    `json:"groups"`
}

// NewNumber adds a new number in database
func (n *Number) NewNumber(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO numbers(userID, phone, verified) VALUES (?, ?, ?)`,
		n.UserID,
		n.Phone,
		n.Verified,
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

// GetNumberByID returns a user given it's ID
func (n *Number) GetNumberByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT userID, phone, verified FROM numbers WHERE id = ?`,
		n.ID,
	).Scan(&n.UserID, &n.Phone, &n.Verified)
}

// GetNumberByPhoneUserID returns a user given it's ID
func (n *Number) GetNumberByPhoneUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, verified FROM numbers WHERE userID = ? AND phone = ?`,
		n.UserID,
		n.Phone,
	).Scan(&n.ID, &n.Verified)
}

func (n *Number) UpdateNumberByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`UPDATE numbers SET userID = ?, phone = ?, verified = ? WHERE id = ?`,
		n.UserID,
		n.Phone,
		n.Verified,
		n.ID,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (n *Number) DeleteNumberByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(`DELETE FROM numbers WHERE id = ?`, n.ID)
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}
