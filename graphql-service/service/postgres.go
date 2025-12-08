package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riyadennis/event-management/graphql-service/internal"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	// Import the database driver (PostgreSQL) and source driver (file system)
	// We use the blank identifier (_) because we only need the init() functions to register the drivers.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func SetUpPostgresDB(ctx context.Context, conf internal.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, getURLForConnectionPool(conf))
	if err != nil {
		return nil, err
	}

	// Test connectivity
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigration(logger *otelzap.Logger, conf internal.Config) error {
	m, err := migrate.New("file://"+conf.MigrationsPath, getDBURLForMigration(conf))

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no migration to run")
		} else {
			logger.Error("failed initialise migration", zap.Error(err))
			return ErrFailedTORunMigration
		}
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("failed to run migration", zap.Error(err))
		return ErrFailedTORunMigration
	}

	return nil
}

func getURLForConnectionPool(conf internal.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)
}

func getDBURLForMigration(conf internal.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
}
