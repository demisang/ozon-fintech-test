package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/demisang/ozon-fintech-test/internal/application"
	"github.com/demisang/ozon-fintech-test/internal/config"
	"github.com/demisang/ozon-fintech-test/internal/rest"
	"github.com/demisang/ozon-fintech-test/internal/store"
	"github.com/demisang/ozon-fintech-test/internal/store/db"
	"github.com/demisang/ozon-fintech-test/internal/store/memory"
	"github.com/demisang/ozon-fintech-test/pkg/grpc"
	"github.com/demisang/ozon-fintech-test/pkg/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	log := logger.GetLogger()

	cfg, err := config.ParseConfig(log)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Init storage
	var storage store.Link

	switch cfg.Storage {
	case "db":
		storage, err = db.New(ctx, log, cfg.Database)

		if err != nil {
			return fmt.Errorf("database storage: %w", err)
		}
	case "memory":
		storage = memory.New(log)
	default:
		return fmt.Errorf("unknown storage '%s'", cfg.Storage)
	}

	// Init app
	app := application.NewApp(log, cfg, storage)

	// Init server
	server := rest.NewServer(log, app, cfg.Server.Host, cfg.Server.Port)

	// Init gRPC server
	grpcServer := grpc.NewLinksServer(app, log, "localhost:9000")

	// Run
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return server.Run(ctx)
	})

	eg.Go(func() error {
		return grpcServer.Run(ctx)
	})

	return eg.Wait()
}
