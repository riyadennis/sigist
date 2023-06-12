package graph

import (
	"database/sql"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var (
	// ErrorFailedToSaveUser means that the user couldn't be saved to db
	ErrorFailedToSaveUser = errors.New("failed to save user")
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// KafkaConfig encapsulates the kafka producer and the topic name
type KafkaConfig struct {
	Topic    string
	Producer Producer
}

type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

// Resolver encapsulates the dependencies for the resolver
type Resolver struct {
	logger      *otelzap.Logger
	db          *sql.DB
	KafkaConfig *KafkaConfig
}

// NewResolver creates a new resolver
func NewResolver(logger *otelzap.Logger, db *sql.DB, kafkaConfig *KafkaConfig) *Resolver {
	return &Resolver{
		logger:      logger,
		db:          db,
		KafkaConfig: kafkaConfig,
	}
}
