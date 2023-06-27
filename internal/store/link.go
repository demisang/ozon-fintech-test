package store

import (
	"context"

	"ozon-fintech-test/internal/models"
)

type Link interface {
	GetByCode(ctx context.Context, code string) (models.Link, error)
	Create(ctx context.Context, createDto models.CreateLinkDto) (models.Link, error)
}
