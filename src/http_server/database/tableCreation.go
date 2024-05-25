package database

import (
	"database/sql"
)

func createTables(db *sql.DB) error {
	err := activateExtension(db)
	if err != nil {
		return err
	}
	err = createUserTable(db)
	if err != nil {
		return err
	}
	err = createPollTable(db)
	if err != nil {
		return err
	}
	err = createPollOptionsTable(db)
	if err != nil {
		return err
	}
	err = createVotesTable(db)
	if err != nil {
		return err
	}
	return nil
}

func createUserTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username CHAR(30) NOT NULL,
		password CHAR(64) NOT NULL
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createPollTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS polls (
		poll_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
		creator_id UUID NOT NULL,
		poll_topic TEXT NOT NULL
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createPollOptionsTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS poll_options (
		poll_id UUID,
		option_num INT,
		option_text TEXT NOT NULL,
		PRIMARY KEY (poll_id, option_num)
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createVotesTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS votes (
		poll_id UUID,
		user_id UUID,
		option_num INT NOT NULL,
		PRIMARY KEY (poll_id, user_id)
	);`
	_, err := db.Exec(sqlStatement)
	return err
}