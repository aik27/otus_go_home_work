package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	log "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"
)

/**
docker compose --env-file configs/calendar.env up -d
goose -dir ./internal/db/migrations postgres "user=otus dbname=calendar
password=otus port=54321 host=localhost sslmode=disable" up
*/

var configFile string

const (
	storageTypeSql    = "sql"
	storageTypeMemory = "memory"
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar.env", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := cfg.New(configFile)
	logger := log.New(config.LogLevel)

	repository, err := getStorage(config.StorageType, config)
	if err != nil {
		panic(err.Error())
	}

	calendar := app.New(logger, config, repository)
	server := internalhttp.NewServer(logger, config, calendar)

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

func getStorage(storageType string, config *cfg.Config) (storage.Repository, error) {
	switch storageType {
	case storageTypeSql:
		res, err := sqlstorage.New(config.GetDbDsn())
		return res, err
	case storageTypeMemory:
		res := memorystorage.New()
		return res, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
