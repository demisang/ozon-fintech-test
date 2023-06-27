package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"ozon-fintech-test/internal/application"
	"ozon-fintech-test/internal/config"
	"ozon-fintech-test/internal/rest"
	"ozon-fintech-test/pkg/logger"
)

const (
	shortLinkLength         = 10
	shortLinkInvalidPattern = `[^a-zA-Z\d_]+`
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

	app, err := application.NewApp(log, cfg)
	if err != nil {
		return fmt.Errorf("new application: %w", err)
	}

	server := rest.NewServer(log, app, cfg.Server.Host, cfg.Server.Port)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return server.Run(ctx)
	})

	return eg.Wait()
}
