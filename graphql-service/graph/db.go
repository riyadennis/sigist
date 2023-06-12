package graph

import (
	"database/sql"
	"github.com/riyadennis/sigist/graphql-service/graph/model"
)

var (
	querySaveUser           = `INSERT INTO user_feedback (id, first_name, last_name, email, job_title, feedback,created_at) VALUES (?, ?, ?, ?, ?, ?,?)`
	queryGetUserByID        = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE id = ?`
	queryGetUserByEmail     = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE email = ?`
	queryGetUserByFirstName = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE first_name = ?`
	queryGetAllUsers        = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback`
)

func saveUserFeedback(db *sql.DB, input model.UserFeedbackInput, uuid, createdAt string) (sql.Result, error) {
	stmt, err := db.Prepare(querySaveUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(
		uuid,
		input.FirstName,
		input.LastName,
		input.Email,
		input.JobTitle,
		input.Feedback,
		createdAt,
	)
}

func getUserRows(db *sql.DB, filter model.FilterInput) (*sql.Rows, error) {
	switch {
	case filter.ID != nil:
		return db.Query(queryGetUserByID, *filter.ID)
	case filter.Email != nil:
		return db.Query(queryGetUserByEmail, filter.Email)
	case filter.FirstName != nil:
		return db.Query(queryGetUserByFirstName, filter.FirstName)
	default:
		return db.Query(queryGetAllUsers)
	}
}
