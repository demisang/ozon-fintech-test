package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"ozon-fintech-test/internal/models"
	"ozon-fintech-test/internal/store"
)

type LinkStorage struct {
	store.Link
	db *DB
}

func NewLinkStorage(db *DB) *LinkStorage {
	return &LinkStorage{
		db: db,
	}
}

func (s *LinkStorage) GetByCode(ctx context.Context, code string) (link models.Link, err error) {
	query := `SELECT url FROM links WHERE code=$1`

	err = s.db.db.QueryRow(ctx, query, code).Scan(&link.URL)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return link, models.ErrLinkNotFound
	case err != nil:
		return link, fmt.Errorf("get link: %w", err)
	}

	link.Code = code

	return link, nil
}

func (s *LinkStorage) Create(ctx context.Context, createDto models.CreateLinkDto) (link models.Link, err error) {
	query := `INSERT INTO links(code, url) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	link.URL = createDto.URL
	link.Code = models.GenerateLinkCodeByURL(createDto.URL)

	_, err = s.db.db.Exec(ctx, query, link.Code, createDto.URL)
	if err != nil {
		return link, fmt.Errorf("insert link error: %w", err)
	}

	return link, nil
}
