package app

import (
	"context"

	cfg "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

func New(logger *logger.Logger, config *cfg.Config, repository storage.Repository) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{Id: id, Title: title})
}

// TODO
