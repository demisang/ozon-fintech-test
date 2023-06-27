package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DB struct {
	log *logrus.Entry
	db  *pgxpool.Pool
}

func New(ctx context.Context, log *logrus.Logger, dsn string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse connection pool cfg: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init connection pool: %w", err)
	}

	db := DB{
		log: log.WithField("module", "store"),
		db:  conn,
	}

	return &db, nil
}
