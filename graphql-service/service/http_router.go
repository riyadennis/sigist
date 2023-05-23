package service

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

func newRouter(srv *handler.Server) http.Handler {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.Recoverer)

	chiRouter.Handle("/", otelhttp.NewHandler(
		playground.Handler("GraphQL playground", "/graphql"),
		"graphql"))

	chiRouter.Handle("/graphql", srv)
	return chiRouter
}
