package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	NEO4j_URI        string `mapstructure:"NEO4j_URI"`
	NEO4jDB_NAME     string `mapstructure:"NEO4jDB_NAME"`
	NEO4jDB_USER     string `mapstructure:"NEO4jDB_USER"`
	NEO4jDB_Password string `mapstructure:"NEO4jDB_Password"`
	PORT             string `mapstructure:"PORT"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			NEO4j_URI:        os.Getenv("NEO4j_URI"),
			NEO4jDB_NAME:     os.Getenv("NEO4jDB_NAME"),
			NEO4jDB_USER:     os.Getenv("NEO4jDB_USER"),
			NEO4jDB_Password: os.Getenv("NEO4jDB_Password"),
			PORT:             os.Getenv("PORT"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.NEO4j_URI == "" {
		err = errors.New("MONGODB_URI is required")
		return
	}

	if config.NEO4jDB_Password == "" {
		err = errors.New("MONGODB_NAME is required")
		return
	}

	if config.NEO4jDB_USER == "" {
		err = errors.New("NEO4jDB_USER is required")
		return
	}
	if config.NEO4jDB_Password == "" {
		err = errors.New("NEO4jDB_Password is required")
		return
	}

	return
}
