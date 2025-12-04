package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/riyadennis/sigist/rest-service/internal"
	"github.com/riyadennis/sigist/rest-service/service"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	server, err := service.NewService(config)
	if err != nil {
		log.Fatal("failed to initialise the service", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("failed to start service", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = server.ShutDown(ctx)
	if err != nil {
		log.Fatal("failed to shut down service", zap.Error(err))
	}
}
