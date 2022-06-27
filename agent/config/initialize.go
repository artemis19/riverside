package config

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

var ConfigFilename string = "agent.yml"

func InitializeConfig() {
	viper.Set("configFile", ConfigFilename)
	viper.SetConfigName(ConfigFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Printf("Config file %v not found. Creating one now.\n", ConfigFilename)

			// Write new config file if one is not found, based on arguments
			viper.SafeWriteConfigAs(ConfigFilename)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Config file %v found but errored with %v", ConfigFilename, err)
		}
	}
	InitializeLogging()
}

func InitializeLogging() {
	// Configure logging by default or supplied arguments
	logFile := viper.Get("outputFile").(string)
	if logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Could not create log file handle.\n")
		}
		if viper.Get("debug").(bool) {
			mw := io.MultiWriter(os.Stdout, file)
			log.SetOutput(mw)
		} else {
			log.SetOutput(file)
		}
	}
}
