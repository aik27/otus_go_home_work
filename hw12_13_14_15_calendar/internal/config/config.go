package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

//nolint:stylecheck
type Config struct {
	Environment string
	Port        int
	StorageType string `split_words:"true"`
	LogLevel    string `split_words:"true"`
	DbHost      string `split_words:"true"`
	DbPort      int    `split_words:"true"`
	DbUsername  string `split_words:"true"`
	DbPassword  string `split_words:"true"`
	DbDatabase  string `split_words:"true"`
}

func (c *Config) GetDbDsn() string {
	//nolint:lll
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow", c.DbHost, c.DbUsername, c.DbPassword, c.DbDatabase, c.DbPort)
}

func New(file string) *Config {
	config := &Config{}

	if err := godotenv.Overload(file); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := envconfig.Process("", config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}

	return config
}
