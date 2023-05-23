package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-migrate/migrate/v4"
	"github.com/riyadennis/sigist/graphql-service/graph"
	"github.com/riyadennis/sigist/graphql-service/graph/generated"
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

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
)

// HTTPServer encapsulates two http server operations  that we need to execute in the service
// it is mainly helpful for testing, by creating mocks for http calls.
type HTTPServer interface {
	Shutdown(ctx context.Context) error
	Serve(l net.Listener) error
}

type Service struct {
	Conf    internal.Config
	Server  HTTPServer
	Logger  *otelzap.Logger
	Sigint  chan os.Signal
	errChan chan error

	DB *sql.DB
}

func NewService(ctx context.Context, conf internal.Config) (*Service, error) {
	log, err := logger(conf.Env)
	if err != nil {
		return nil, err
	}

	logger := otelzap.New(log)

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		logger.Error("failed to open db connection", zap.Error(err))
		return nil, ErrFailedTOOpenDB
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: graph.NewResolver(
					logger,
					db,
				),
			}),
	)
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		logger.Error("failed initialise Db driver", zap.Error(err))
		return nil, ErrFailedTOOpenDB
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no migration to run")
		} else {
			logger.Error("failed initialise migration", zap.Error(err))
			return nil, ErrFailedTORunMigration
		}
	}

	err = m.Up()
	if err != migrate.ErrNoChange && err != nil {
		logger.Error("failed to run migration", zap.Error(err))
		return nil, ErrFailedTORunMigration
	}
	server := &http.Server{
		Addr:    conf.Port,
		Handler: newRouter(srv),
	}
	return &Service{
		Conf:   conf,
		Logger: logger,
		Server: server,
		DB:     db,
	}, nil
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

func logger(env string) (*zap.Logger, error) {
	log := &zap.Logger{}
	var err error
	switch env {
	case "prod":
		log, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case "test":
		log = zap.NewExample()
	case "dev":
		log, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	return log, nil
}
