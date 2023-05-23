package graph

import (
	"database/sql"
	"github.com/riyadennis/sigist/graphql-service/graph/model"
)

var (
	querySaveUser       = `INSERT INTO users (first_name, last_name, email, job_title, created_at) VALUES (?, ?, ?, ?, ?)`
	queryGetUserByID    = `SELECT * FROM users WHERE id = ?`
	queryGetUserByEmail = `SELECT * FROM users WHERE email = ?`
	queryGetAllUsers    = `SELECT * FROM users`
)

func saveUser(db *sql.DB, input model.CreateUserInput, createdAt string) (sql.Result, error) {
	stmt, err := db.Prepare(querySaveUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(
		input.FirstName,
		input.LastName,
		input.Email,
		input.JobTitle,
		createdAt,
	)
}

func getUserRows(db *sql.DB, filter model.FilterInput) (*sql.Rows, error) {
	switch {
	case filter.ID != nil:
		return db.Query(queryGetUserByID, *filter.ID)
	case filter.Email != nil:
		return db.Query(queryGetUserByEmail, filter.Email)
	default:
		return db.Query(queryGetAllUsers)
	}
}
