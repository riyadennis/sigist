package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	// Import the database driver (PostgreSQL) and source driver (file system)
	// We use the blank identifier (_) because we only need the init() functions to register the drivers.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	// ErrFailedTORunMigration means that the migration couldn't be run
	ErrFailedTORunMigration = errors.New("failed to run migration")
)

type Config struct {
	Host           string `arg:"env:DB_HOST" default:"localhost"`
	User           string `arg:"env:DB_USER" default:"username"`
	Password       string `arg:"env:DB_PASSWORD" default:"password"`
	Name           string `arg:"env:DB_NAME" default:"feedback"`
	Port           string `arg:"env:DB_PORT" default:"5434"`
	MigrationsPath string `arg:"env:MIGRATIONS_PATH" default:"migrations"`
}

func NewConfig() (Config, error) {
	var conf Config
	validate := validator.New()
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return conf, fmt.Errorf("failed to register \"NotBlank\" validator")
	}

	arg.MustParse(&conf)
	if err := validate.Struct(conf); err != nil {
		return conf, fmt.Errorf("validating struct: %w", err)
	}

	return conf, err
}

func setUpPostgresDB(ctx context.Context, conf Config) (*pgxpool.Pool, error) {
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

func runMigration(logger *otelzap.Logger, conf Config) error {
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

func getURLForConnectionPool(conf Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
}

func getDBURLForMigration(conf Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
}
