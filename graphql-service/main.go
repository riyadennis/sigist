package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/riyadennis/sigist/graphql-service/service"
	"go.uber.org/zap"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	server, err := service.NewService(config)
	if err != nil {
		log.Fatal("failed to initialise service ", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("failed to start service", err)
	}

	ctx := context.Background()
	err = server.ShutDown(ctx)
	if err != nil {
		log.Fatal("failed to shut down service", zap.Error(err))
	}
}
