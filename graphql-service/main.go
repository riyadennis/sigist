package main

import (
	"context"
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/riyadennis/sigist/graphql-service/service"
	"go.uber.org/zap"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	ctx := context.Background()
	server, err := service.NewService(ctx, config)
	if err != nil {
		log.Fatal("failed to initialise service", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("failed to start service", err)
	}

	err = server.ShutDown(ctx)
	if err != nil {
		log.Fatal("failed to shut down service", zap.Error(err))
	}
}
