package main

import (
	"project/pkg/db/cache"
	"project/pkg/db/sqlx"
	"project/pkg/server"

	"github.com/spf13/viper"
)

type CoreConfig struct {
	Server *server.Config
	DB     *sqlx.DbConfig
	Cache  cache.Config
}

var config CoreConfig

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("core")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.Unmarshal(&config)
	return nil
}

func NewCacheConfig() *cache.Config {
	return &config.Cache
}
