package memory

import (
	"context"
	"errors"
	"fmt"

	"ozon-fintech-test/internal/models"
	"ozon-fintech-test/internal/store"
)

type LinkStorage struct {
	store.Link
	memory *Memory
}

func NewLinkStorage(memory *Memory) *LinkStorage {
	return &LinkStorage{
		memory: memory,
	}
}

func (s *LinkStorage) GetByCode(ctx context.Context, code string) (link models.Link, err error) {
	url, err := s.memory.get(ctx, code)

	switch {
	case errors.Is(err, ErrNotExists):
		return link, models.ErrLinkNotFound
	case err != nil:
		return link, fmt.Errorf("get link: %w", err)
	}

	link.Code = code
	link.URL = url

	return link, nil
}

func (s *LinkStorage) Create(ctx context.Context, createDto models.CreateLinkDto) (models.Link, error) {
	code := models.GenerateLinkCodeByURL(createDto.URL)
	s.memory.set(ctx, code, createDto.URL)

	return models.Link{
		Code: code,
		URL:  createDto.URL,
	}, nil
}
