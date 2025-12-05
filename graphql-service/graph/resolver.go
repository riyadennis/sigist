package graph

import (
	"context"
	"database/sql"
	"errors"

	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var (
	// ErrorFailedToSaveUser means that the user couldn't be saved to db
	ErrorFailedToSaveUser = errors.New("failed to save user")
)

// Resolver encapsulates the dependencies for the resolver
type Resolver struct {
	logger      *otelzap.Logger
	db          *sql.DB
	KafkaConfig *internal.KafkaConfig
	KafkaWriter *internal.KafkaWriter
}

// NewResolver creates a new resolver
func NewResolver(ctx context.Context, logger *otelzap.Logger, db *sql.DB, kafkaConfig *internal.KafkaConfig) (*Resolver, error) {
	connection, err := kafkaConfig.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return &Resolver{
		logger:      logger,
		db:          db,
		KafkaConfig: kafkaConfig,
		KafkaWriter: &internal.KafkaWriter{Connection: connection},
	}, nil
}
