package application

import (
	"github.com/demisang/ozon-fintech-test/internal/config"
	"github.com/demisang/ozon-fintech-test/internal/store"
	"github.com/sirupsen/logrus"
)

type App struct {
	log     *logrus.Entry
	Storage store.Link
	Config  config.Config
}

func NewApp(log *logrus.Logger, cfg config.Config, storage store.Link) (*App, error) {
	a := App{
		log:     log.WithField("module", "application"),
		Storage: storage,
		Config:  cfg,
	}

	return &a, nil
}
