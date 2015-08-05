package main

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

func loadConfig(filenamePath *string, filename *string) {
	log := logging.MustGetLogger("log")

	viper.SetConfigName(*filename)
	viper.AddConfigPath(*filenamePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Critical("Unable to load config file:", err)
	}
}
