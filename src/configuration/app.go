package configuration

import (
	"github.com/marcellinuskristanto/onepiece/src/helper"
	"github.com/spf13/viper"
)

// AppConfiguration type
type AppConfiguration struct {
	Listen int
	Env    string
}

func loadAppConfiguration() (AppConfiguration, error) {
	provider := viper.New()

	provider.SetConfigName("app")
	provider.AddConfigPath("./config")

	var config AppConfiguration

	setDefaultAppConfiguration(provider)

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

func setDefaultAppConfiguration(provider *viper.Viper) {
	provider.SetDefault("Listen", helper.GetEnv("LISTEN", 3000).(int))
	provider.SetDefault("Env", helper.GetEnv("ENV", "production").(string))
}
