package service

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/riyadennis/sigist/rest-service/internal"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// these two can be overridden at build time
	serviceVersion    = "DEV"
	serviceCommitHash = ""

	// ErrFailedToStartListener means that the listener couldn't be started
	ErrFailedToStartListener = errors.New("failed to start listener")

	// ErrFailedToStartServer means that the server couldn't be started
	ErrFailedToStartServer = errors.New("failed to start server")

	// ErrFailedTOOpenDB means that the db couldn't be opened
	ErrFailedTOOpenDB = errors.New("failed to open db")

	// ErrFailedTORunMigration means that the migration couldn't be run
	ErrFailedTORunMigration = errors.New("failed to run migration")

	// ErrFailedToCreateKafkaProducer means that the kafka producer couldn't be created
	ErrFailedToCreateKafkaProducer = errors.New("failed to create kafka producer")
)

// HTTPServer encapsulates two http server operations  that we need to execute in the service
// it is mainly helpful for testing, by creating mocks for http calls.
type HTTPServer interface {
	Shutdown(ctx context.Context) error
	Serve(l net.Listener) error
}

// Service encapsulates the service operations
type Service struct {
	Conf    internal.Config
	Server  HTTPServer
	Logger  *otelzap.Logger
	Sigint  chan os.Signal
	errChan chan error
	DB      *sql.DB
}

// NewService creates a new service
func NewService(conf internal.Config) (*Service, error) {
	log, err := logger(conf.Env)
	if err != nil {
		return nil, err
	}

	logger := otelzap.New(log)
	db, err := SetUpDB(conf.DBFile, conf.MigrationsPath)
	if err != nil {
		logger.Error("failed to open db connection", zap.Error(err))
		return nil, ErrFailedTOOpenDB
	}
	server := &http.Server{
		Addr:    conf.Port,
		Handler: newRouter(db, logger),
	}

	return &Service{
		Conf:   conf,
		Logger: logger,
		Server: server,
		DB:     db,
	}, nil
}

func SetUpDB(dbFile, migrationsPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3", driver)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, ErrFailedTORunMigration
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, ErrFailedTORunMigration
	}

	return db, nil
}

// Start the service will kick-start http server, kafka and other needed processes.
// If an error is returned then the http listener goroutine has been started.
func (s *Service) Start() error {
	s.Logger.Info("starting service",
		zap.String("port", s.Conf.Port),
		zap.String("version", serviceVersion),
		zap.String("commit hash", serviceCommitHash),
	)
	listener, err := net.Listen("tcp", s.Conf.Port)
	if err != nil {
		s.Logger.Error("failed to start http listener", zap.Error(err))
		return ErrFailedToStartListener
	}

	go func() {
		s.Logger.Info("service finished starting and is now ready to accept requests")

		// start http listener
		err := s.Server.Serve(listener)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.Logger.Error("failed to start http server", zap.Error(err))
				s.errChan <- ErrFailedToStartServer
				return
			}
		}
	}()

	return nil
}

// ShutDown will wait for error in error channel or signal interrupt in signal channel
func (s *Service) ShutDown(ctx context.Context) error {
	select {
	case err := <-s.errChan:
		if err != nil {
			return err
		}
	case <-s.Sigint:
		close(s.errChan)
		s.gracefulShutdown(ctx)
		s.DB.Close()
		return nil
	}

	return nil
}

// gracefulShutdown gracefully shutdown the service and its dependencies
func (s *Service) gracefulShutdown(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		_ = s.Logger.Sync()
		cancel()
	}()

	_ = s.Server.Shutdown(cancelCtx)
}

func newRouter(db *sql.DB, logger *otelzap.Logger) http.Handler {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	eh := NewEmailHandler(db, logger)
	chiRouter.MethodFunc(http.MethodPost, "/email", eh.SaveEmail)
	chiRouter.MethodFunc(http.MethodGet, "/emails", eh.GetAllEmails)
	return chiRouter
}

func logger(env string) (*zap.Logger, error) {
	switch env {
	case "prod":
		return zap.NewProduction()
	case "test":
		return zap.NewExample(), nil
	case "dev":
		return zap.NewDevelopment()
	default:
		return nil, errors.New("invalid environment")
	}
}
