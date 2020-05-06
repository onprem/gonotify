package models

import "database/sql"

// Group represents a group of nodes
type Group struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

// NewGroup creates a new group in database
func (g *Group) NewGroup(tx *sql.Tx) (int64, error) {
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

	return id, nil
}

// GetGroupByID returns a group given it's ID
func (g *Group) GetGroupByID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT name, userID FROM groups WHERE id = ?`,
		g.ID,
	).Scan(&g.Name, &g.UserID)
}

// GetGroupByUserID returns a group given it's userID
func (g *Group) GetGroupByUserID(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, name FROM groups WHERE userID = ?`,
		g.UserID,
	).Scan(&g.ID, &g.Name)
}
