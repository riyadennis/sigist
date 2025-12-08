package graph

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riyadennis/event-management/graphql-service/graph/model"
)

var (
	querySaveUser           = `INSERT INTO user_feedback (id,first_name,last_name, email,job_title,feedback) VALUES ($1, $2, $3, $4, $5, $6)`
	queryGetUserByID        = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE id = $1`
	queryGetUserByEmail     = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE email = $1`
	queryGetUserByFirstName = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback WHERE first_name = $1`
	queryGetAllUsers        = `SELECT id, first_name, last_name, email, job_title, feedback, created_at FROM user_feedback`
)

func saveUserFeedback(ctx context.Context, pool *pgxpool.Pool, input model.UserFeedbackInput, uuid string) (int, error) {
	stmt, err := pool.Exec(
		ctx,
		querySaveUser,
		uuid,
		input.FirstName,
		input.LastName,
		input.Email,
		input.JobTitle,
		input.Feedback)
	if err != nil {
		return 0, err
	}
	return int(stmt.RowsAffected()), nil
}

func getUserRows(ctx context.Context, pool *pgxpool.Pool, filter model.FilterInput) (pgx.Rows, error) {
	switch {
	case filter.ID != nil:
		return pool.Query(ctx, queryGetUserByID, *filter.ID)
	case filter.Email != nil:
		return pool.Query(ctx, queryGetUserByEmail, filter.Email)
	case filter.FirstName != nil:
		return pool.Query(ctx, queryGetUserByFirstName, filter.FirstName)
	default:
		return pool.Query(ctx, queryGetAllUsers)
	}
}
