package application

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"ozon-fintech-test/internal/config"
	"ozon-fintech-test/internal/store"
	"ozon-fintech-test/internal/store/db"
	"ozon-fintech-test/internal/store/memory"
)

type App struct {
	log     *logrus.Entry
	Storage store.Link
	Config  config.Config
}

func NewApp(log *logrus.Logger, cfg config.Config) (*App, error) {
	var storage store.Link
	var err error

	ctx := context.Background()
	switch cfg.Storage {
	case "db":
		storage, err = NewDBStorage(ctx, log, cfg.Database)
	case "memory":
		storage, err = NewMemoryStorage(ctx, log)
	default:
		return nil, fmt.Errorf("unknown storage '%s'", cfg.Storage)
	}

	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	a := App{
		log:     log.WithField("module", "application"),
		Storage: storage,
		Config:  cfg,
	}

	return &a, nil
}

func NewDBStorage(ctx context.Context, log *logrus.Logger, dbConfig config.Database) (*db.LinkStorage, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbConfig.Username, dbConfig.Password),
		Host:     dbConfig.Host + ":" + dbConfig.Port,
		Path:     dbConfig.DBName,
		RawQuery: "sslmode=disable",
	}

	database, err := db.New(ctx, log, dsn.String())
	if err != nil {
		return nil, fmt.Errorf("database storage: %w", err)
	}

	return db.NewLinkStorage(database), nil
}

func NewMemoryStorage(ctx context.Context, log *logrus.Logger) (*memory.LinkStorage, error) {
	mem, err := memory.New(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("memory storage: %w", err)
	}

	return memory.NewLinkStorage(mem), nil
}
