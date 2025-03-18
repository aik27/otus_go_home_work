package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	app    Application
	server *http.Server
}

type Application interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func New(config *config.Config, logger *logger.Logger, app Application) *Server {
	return &Server{
		config: config,
		logger: logger,
		app:    app,
		server: &http.Server{
			Addr:         net.JoinHostPort(config.HttpServerHost, config.HttpServerPort),
			ReadTimeout:  time.Duration(config.HttpServerReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.HttpServerWriteTimeout) * time.Second,
		},
	}
}

type HelloHandler struct{}

func (h *HelloHandler) Hello(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello world!"))
	if err != nil {
		return
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.app.Start(ctx)
	if err != nil {
		return err
	}

	handler := &HelloHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingMiddleware(s.logger, handler.Hello))

	s.server.Handler = mux

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.app.Stop(ctx)
	if err != nil {
		return err
	}
	return nil
}
