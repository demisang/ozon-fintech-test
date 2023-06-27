package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/demisang/ozon-fintech-test/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DB struct {
	log *logrus.Entry
	db  *pgxpool.Pool
}

func New(ctx context.Context, log *logrus.Logger, dbConfig config.Database) (*LinkStorage, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbConfig.Username, dbConfig.Password),
		Host:     dbConfig.Host + ":" + dbConfig.Port,
		Path:     dbConfig.DBName,
		RawQuery: "sslmode=disable",
	}
	cfg, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, fmt.Errorf("pgsql error: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init connection pool: %w", err)
	}

	db := DB{
		log: log.WithField("module", "store"),
		db:  conn,
	}

	return newLinkStorage(&db), nil
}
