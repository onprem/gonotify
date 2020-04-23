package api

import (
	"database/sql"
)

func bootstrapDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		name TEXT,
		phone TEXT UNIQUE,
		hash TEXT,
		verified BOOLEAN NOT NULL CHECK (verified IN (0,1)))`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		userID INTEGER NOT NULL,
		name TEXT)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS numbers(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		groupID INTEGER NOT NULL,
		phone TEXT UNIQUE,
		verified BOOLEAN NOT NULL CHECK (verified IN (0,1)),
		lastMsgReceived TEXT)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pendingMsgs(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		notifID INTEGER NOT NULL,
		numberID INTEGER NOT NULL)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notifications(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		userID INTEGER NOT NULL,
		groupID INTEGER NOT NULL,
		body TEXT,
		timeSt TEXT)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS numberVerify(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		numberID INTEGER NOT NULL,
		code TEXT)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS userVerify(
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		userID INTEGER NOT NULL,
		numberID INTEGER NOT NULL,
		code TEXT)`,
	)
	if err != nil {
		return err
	}

	return nil
}
