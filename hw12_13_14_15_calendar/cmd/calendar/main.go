package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	log "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := cfg.NewConfig(configFile)
	logger := log.New(config.LogLevel)

	storage := memorystorage.New()
	calendar := app.New(logger, storage)
	server := internalhttp.NewServer(logger, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
