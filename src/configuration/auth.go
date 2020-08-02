package configuration

import (
	"github.com/marcellinuskristanto/onepiece/src/helper"
	"github.com/spf13/viper"
)

// AuthConfiguration type
type AuthConfiguration struct {
	Secret   string
	Username string
	Password string
}

func loadAuthConfiguration() (AuthConfiguration, error) {
	provider := viper.New()

	provider.SetConfigName("auth")
	provider.AddConfigPath("./config")

	var config AuthConfiguration

	setDefaultAuthConfiguration(provider)

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

func setDefaultAuthConfiguration(provider *viper.Viper) {
	provider.SetDefault("Secret", helper.GetEnv("SECRET", "greenlandluffy").(string))
	provider.SetDefault("Username", helper.GetEnv("USERNAME", "randomusernamesolong").(string))
	provider.SetDefault("Password", helper.GetEnv("PASSWORD", "RndomP4dsssword").(string))
}
