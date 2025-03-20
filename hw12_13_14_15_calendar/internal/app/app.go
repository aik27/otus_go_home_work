package app

import (
	"context"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	config  *config.Config
	logger  *logger.Logger
	storage storage.Repository
}

func New(config *config.Config, logger *logger.Logger, storage storage.Repository) *App {
	return &App{
		config:  config,
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Start(ctx context.Context) error {
	return a.storage.Connect(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	return a.storage.Close(ctx)
}

func (a *App) CreateEvent(e storage.Event) error {
	if _, err := a.storage.Save(e); err != nil {
		a.logger.Error("create event error: " + err.Error())
		return err
	}
	return nil
}
