package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	//"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Port        int
	LogLevel    string `split_words:"true"`
	DbUsername  string `split_words:"true"`
	DbPassword  string `split_words:"true"`
	DbDatabase  string `split_words:"true"`
}

func NewConfig(file string) Config {
	config := Config{}

	if err := godotenv.Overload(file); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := envconfig.Process("", &config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}

	/*
		viper.SetConfigFile(file)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %w", err))
		}
	*/
	return config
}
