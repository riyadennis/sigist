package main

import (
	"context"
	"log"

	"github.com/riyadennis/event-management/graphql-service/internal"
	"github.com/riyadennis/event-management/graphql-service/service"
	"go.uber.org/zap"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := service.NewService(ctx, config)
	if err != nil {
		log.Fatal("failed to initialise service ", err)
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
