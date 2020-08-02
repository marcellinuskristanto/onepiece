package configuration

// Configuration type
type Configuration struct {
	App    AppConfiguration
	Auth   AuthConfiguration
	Logger LoggerConfiguration
}

// LoadConfigurations load config file
func LoadConfigurations() (config Configuration, err error) {
	appConfig, err := loadAppConfiguration()
	if err != nil {
		return config, err
	}
	config.App = appConfig

	loggerConfig, err := loadLoggerConfiguration()
	if err != nil {
		return config, err
	}
	config.Logger = loggerConfig

	authConfig, err := loadAuthConfiguration()
	if err != nil {
		return config, err
	}
	config.Auth = authConfig

	return config, nil
}
