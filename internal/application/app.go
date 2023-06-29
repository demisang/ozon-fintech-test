package application

import (
	"context"
	"fmt"

	"github.com/demisang/ozon-fintech-test/internal/config"
	"github.com/demisang/ozon-fintech-test/internal/models"
	"github.com/demisang/ozon-fintech-test/internal/store"
	"github.com/sirupsen/logrus"
)

type App struct {
	log     *logrus.Entry
	storage store.Link
	Config  config.Config
}

func NewApp(log *logrus.Logger, cfg config.Config, storage store.Link) *App {
	return &App{
		log:     log.WithField("module", "application"),
		storage: storage,
		Config:  cfg,
	}
}

func (a *App) LinkGet(ctx context.Context, code string) (models.Link, error) {
	link, err := a.storage.GetByCode(ctx, code)
	if err != nil {
		return link, fmt.Errorf("app get link: %w", err)
	}

	return link, nil
}

func (a *App) LinkCreate(ctx context.Context, createDto models.CreateLinkDTO) (models.Link, error) {
	link, err := a.storage.Create(ctx, createDto)
	if err != nil {
		return link, fmt.Errorf("app create link: %w", err)
	}

	return link, nil
}

func (a *App) ValidateLinkCodeLength(code string) bool {
	return len(code) == a.Config.ShortLinkLength
}
