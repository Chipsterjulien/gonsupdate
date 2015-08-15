package main

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"
)

func loadConfig(filenamePath *string, filename *string) {
	log := logging.MustGetLogger("log")

	viper.SetConfigName(*filename)
	viper.AddConfigPath(*filenamePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Critical("Unable to load config file:", err)
		os.Exit(1)
	}

	switch viper.GetString("logtype") {
	case "critical":
		logging.SetLevel(0, "")
		log.Debug("\t\"critical\" is selected")
	case "error":
		logging.SetLevel(1, "")
		log.Debug("\t\"error\" is selected")
	case "warning":
		logging.SetLevel(2, "")
		log.Debug("\t\"warning\" is selected")
	case "notice":
		logging.SetLevel(3, "")
		log.Debug("\t\"notice\" is selected")
	case "info":
		logging.SetLevel(4, "")
		log.Debug("\t\"info\" is selected")
	case "debug":
		logging.SetLevel(5, "")
		log.Debug("\t\"debug\" is selected")
	default:
		logging.SetLevel(2, "")
		log.Debug("\t\"default\" is selected (warning)")
	}

	log.Debug("loadConfig func:")
	log.Debug("  path: %s", *filenamePath)
	log.Debug("  filename: %s", *filename)
	log.Debug("  logtype in file config is \"%s\"", viper.GetString("logtype"))
}
