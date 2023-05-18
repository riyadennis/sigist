package graph

import (
	"database/sql"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger *otelzap.Logger
	db     *sql.DB
}

func NewResolver(logger *otelzap.Logger, db *sql.DB) *Resolver {
	return &Resolver{
		logger: logger,
		db:     db,
	}
}
