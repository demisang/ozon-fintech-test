package store

import (
	"context"

	"github.com/demisang/ozon-fintech-test/internal/models"
)

type Link interface {
	GetByCode(ctx context.Context, code string) (models.Link, error)
	Create(ctx context.Context, createDto models.CreateLinkDTO) (models.Link, error)
}
