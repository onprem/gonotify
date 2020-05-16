package models

import "database/sql"

// Number represnts a phone number
type Number struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userID"`
	Phone    string `json:"phone"`
	Verified bool   `json:"verified"`
	Groups   int    `json:"groups"`
}

// New adds a new number in database
func (n *Number) New(tx *sql.Tx) (int64, error) {
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
	n.ID = int(id)
	return id, nil
}

// GetByID returns a user given it's ID
func (n *Number) GetByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT userID, phone, verified FROM numbers WHERE id = ?`,
		n.ID,
	).Scan(&n.UserID, &n.Phone, &n.Verified)
}

// GetByPhoneUserID returns a user given it's ID
func (n *Number) GetByPhoneUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, verified FROM numbers WHERE userID = ? AND phone = ?`,
		n.UserID,
		n.Phone,
	).Scan(&n.ID, &n.Verified)
}

// UpdateByID updates a number given it's ID
func (n *Number) UpdateByID(tx *sql.Tx) (int64, error) {
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

// DeleteByID deletes a number given it's ID
func (n *Number) DeleteByID(tx *sql.Tx) (int64, error) {
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
