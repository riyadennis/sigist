package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var (
	querySaveEmail    = `INSERT INTO emails (id, sourceName, email, created_at) VALUES (?, ?, ?, ?)`
	queryGetAllEmails = `SELECT id, sourceName, email, created_at FROM emails`
)

type Email struct {
	logger *otelzap.Logger
	db     *sql.DB
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

func NewEmailHandler(db *sql.DB, logger *otelzap.Logger) *Email {
	return &Email{
		logger: logger,
		db:     db,
	}
}

func (e *Email) SaveEmail(w http.ResponseWriter, r *http.Request) {
	request := json.NewDecoder(r.Body)
	re := &Request{}
	err := request.Decode(&re)
	if err != nil {
		_ = HTTPResponse(w, err, http.StatusBadRequest, "failed to decode request")
		return
	}
	emailID := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	res, err := SaveEmail(e.db, re, emailID, createdAt)
	if err != nil {
		e.logger.Error("failed to execute statement", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusInternalServerError, "failed to save email")
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		e.logger.Error("failed to fetch result from db after saving email", zap.Error(err))
		_ = HTTPResponse(w, err, http.StatusInternalServerError, "failed to fetch result from db after saving email")
		return
	}

	if rows == 0 {
		e.logger.Debug("No email saved")
		_ = HTTPResponse(w, nil, http.StatusOK, "success")
		return
	}

	_ = HTTPResponse(w, nil, http.StatusCreated, "success")
}

func (e *Email) GetAllEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := FetchEmails(e.db, e.logger)
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

func FetchEmails(db *sql.DB, logger *otelzap.Logger) ([]*EmailResponse, error) {
	rows, err := db.Query(queryGetAllEmails)
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

func SaveEmail(db *sql.DB, req *Request, uuid, createdAt string) (sql.Result, error) {
	stmt, err := db.Prepare(querySaveEmail)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if len(req.Sources) == 0 {
		req.Sources = []string{"default"}
	}
	return stmt.Exec(
		uuid,
		strings.Join(req.Sources, ","),
		req.Email,
		createdAt,
	)
}
