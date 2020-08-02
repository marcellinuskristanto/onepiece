package configuration

import (
	"github.com/marcellinuskristanto/onepiece/src/helper"
	"github.com/spf13/viper"
)

// LoggerConfiguration type
type LoggerConfiguration struct {
	Path string
}

func loadLoggerConfiguration() (LoggerConfiguration, error) {
	provider := viper.New()

	provider.SetConfigName("logger")
	provider.AddConfigPath("./config")

	var config LoggerConfiguration

	setDefaultLoggerConfiguration(provider)

	err := provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return config, err
		}
	}

	err = provider.Unmarshal(&config)

	return config, err
}

func setDefaultLoggerConfiguration(provider *viper.Viper) {
	provider.SetDefault("Path", helper.GetEnv("LOG_PATH", ""))
}
