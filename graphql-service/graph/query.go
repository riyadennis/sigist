package graph

import (
	"database/sql"
)

var (
	queryFetchID      = `SELECT id FROM users WHERE email = ?`
	queryFetchDetails = `SELECT first_name,first_name,email,job_title FROM users WHERE id = ?`
)

func fetchID(db *sql.DB, query, argument string) (string, error) {
	var result string
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	rows, err := stmt.Query(argument)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
