package internal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/riyadennis/event-management/graphql-service/graph/model"
)

type Producer interface {
	Write(msg *model.UserFeedback) error
}

type KafkaConnection interface {
	Write([]byte) (int, error)
	Close() error
}

type KafkaWriter struct {
	Connection KafkaConnection
}

type KafkaConfig struct {
	Topic     string
	Partition int
	Address   string
	Writer    *KafkaWriter
}

// KafkaSetup creates the config object and opens a kafka connection
func KafkaSetup(ctx context.Context, conf Config) (*KafkaConfig, error) {
	kc := &KafkaConfig{
		Topic:     conf.KafkaTopic,
		Partition: conf.KafkaPartition,
		Address:   conf.KafkaBroker,
	}
	conn, err := kc.Connection(ctx)
	if err != nil {
		return nil, err
	}
	kc.Writer = &KafkaWriter{Connection: conn}

	return kc, nil
}

func (kc *KafkaConfig) Connection(ctx context.Context) (KafkaConnection, error) {
	conn, err := kafka.DialLeader(ctx, "tcp", kc.Address, kc.Topic, kc.Partition)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now().Add(120 * time.Second))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (kc *KafkaWriter) Write(msg *model.UserFeedback) error {
	kafkaData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = kc.Connection.Write(kafkaData)
	if err != nil {
		return err
	}

	return nil
}
