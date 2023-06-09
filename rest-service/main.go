package main

import (
	"context"
	"go.uber.org/zap"
	"log"

	"github.com/riyadennis/sigist/rest-service/internal"
	"github.com/riyadennis/sigist/rest-service/service"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	server, err := service.NewService(config)
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
