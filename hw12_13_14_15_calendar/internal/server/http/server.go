package internalhttp

import (
	"context"
	cfg "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

type Server struct { // TODO
}

type Application interface { // TODO
}

func NewServer(logger *logger.Logger, app Application, config *cfg.Config) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
