package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var (
	queryGetAllEmails = `SELECT id, sourceNames, email, created_at FROM emails`
)

type Email struct {
	logger *otelzap.Logger
	pool   *pgxpool.Pool
}

type Request struct {
	Email   string   `json:"email"`
	Sources []string `json:"sources"`
}
type EmailResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	SourceName string `json:"source_name"`
	CreatedAt  string `json:"created_at"`
}

func NewEmailHandler(db *pgxpool.Pool, logger *otelzap.Logger) *Email {
	return &Email{
		logger: logger,
		pool:   db,
	}
}

func (e *Email) SaveEmail(w http.ResponseWriter, r *http.Request) {
	request := json.NewDecoder(r.Body)
	re := &Request{}
	err := request.Decode(&re)
	if err != nil {
		e.logger.Error("failed to decode request, got error", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusBadRequest, "failed to decode request")
		return
	}
	emailID := uuid.New().String()
	res, err := SaveEmail(r.Context(), e.pool, re, emailID)
	if err != nil {
		e.logger.Error("failed to execute statement", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusInternalServerError, "failed to save email")
		return
	}

	if res == 0 {
		e.logger.Debug("No email saved")
		_ = HTTPResponse(w, nil, http.StatusOK, "success")
		return
	}

	_ = HTTPResponse(w, nil, http.StatusCreated, "success")
}

func (e *Email) GetAllEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := FetchEmails(r.Context(), e.pool, e.logger)
	if err != nil {
		e.logger.Error("failed to fetch emails", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusInternalServerError, "failed to fetch emails")
		return
	}

	data, err := json.Marshal(emails)
	if err != nil {
		e.logger.Error("failed to marshal emails", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusInternalServerError, "failed to fetch emails")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func FetchEmails(ctx context.Context, pool *pgxpool.Pool, logger *otelzap.Logger) ([]*EmailResponse, error) {
	rows, err := pool.Query(ctx, queryGetAllEmails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	emails := make([]*EmailResponse, 0)
	for rows.Next() {
		response := &EmailResponse{}
		err := rows.Scan(&response.ID, &response.SourceName, &response.Email, &response.CreatedAt)
		if err != nil {
			logger.Error("failed to scan row", zap.Error(err))
			return nil, err
		}
		emails = append(emails, response)
		logger.Debug("email", zap.String("id", response.ID),
			zap.String("source", response.SourceName),
			zap.String("email", response.Email),
			zap.String("created_at", response.CreatedAt),
		)

	}

	return emails, nil
}

var querySaveEmail = `INSERT INTO emails (id, sourceNames, email) VALUES ($1, $2, $3)`

func SaveEmail(ctx context.Context, pool *pgxpool.Pool, req *Request, uuid string) (int, error) {
	if len(req.Sources) == 0 {
		req.Sources = []string{"default"}
	}

	cmd, err := pool.Exec(ctx, querySaveEmail, uuid, req.Sources, req.Email)
	if err != nil {
		return 0, err
	}

	return int(cmd.RowsAffected()), nil
}
