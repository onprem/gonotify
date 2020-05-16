package models

import "database/sql"

// Group represents a group of nodes
type Group struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

// New creates a new group in database
func (g *Group) New(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`INSERT INTO groups(name, userID) VALUES (?, ?)`,
		g.Name,
		g.UserID,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	g.ID = int(id)
	return id, nil
}

// GetByID returns a group given it's ID
func (g *Group) GetByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT name, userID FROM groups WHERE id = ?`,
		g.ID,
	).Scan(&g.Name, &g.UserID)
}

// GetByNameUserID returns a group given it's userID
func (g *Group) GetByNameUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, name FROM groups WHERE name = ? AND userID = ?`,
		g.Name,
		g.UserID,
	).Scan(&g.ID, &g.Name)
}

// DeleteByID delets the Group using given ID
func (g *Group) DeleteByID(tx *sql.Tx) (int64, error) {
	res, err := tx.Exec(
		`DELETE FROM groups WHERE id = ?`,
		g.ID,
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
