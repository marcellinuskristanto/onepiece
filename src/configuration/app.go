package configuration

import (
	"github.com/marcellinuskristanto/onepiece/src/helper"
	"github.com/spf13/viper"
)

// AppConfiguration type
type AppConfiguration struct {
	Port        int
	Env         string
	MinioUrl    string
	MinioUser   string
	MinioSecret string
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
	provider.SetDefault("Port", helper.GetEnvInt("PORT", 80))
	provider.SetDefault("Env", helper.GetEnv("ENV", "production"))
	provider.SetDefault("MinioUrl", helper.GetEnv("MINIO_URL", "play.min.io"))
	provider.SetDefault("MinioUser", helper.GetEnv("MINIO_USER", "user"))
	provider.SetDefault("MinioSecret", helper.GetEnv("MINIO_SECRET", "secret"))
}
