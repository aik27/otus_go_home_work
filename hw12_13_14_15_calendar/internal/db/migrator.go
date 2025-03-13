package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"log"
)

func Migrate(db *sqlx.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose: failed to set dialect: %v", err)
	}

	if err := goose.Up(db.DB, "internal/db/migrations"); err != nil {
		log.Fatalf("goose up error: %v", err)
	}
}
