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
	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/db"
	log "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var configFile string

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

	config := cfg.NewConfig(configFile)
	logger := log.New(config.LogLevel)

	database, err := sqlx.Open("pgx", config.GetDbDsn())
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close database connection: %w", err))
		}
	}(database)

	if err != nil {
		panic(fmt.Errorf("failed to load driver: %w", err))
	}

	err = database.PingContext(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	db.Migrate(database)

	storage := memorystorage.New()
	calendar := app.New(logger, storage, config)
	server := internalhttp.NewServer(logger, calendar, config)

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
