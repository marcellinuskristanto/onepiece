package configuration

// Configuration type
type Configuration struct {
	App    AppConfiguration
	Auth   AuthConfiguration
	Logger LoggerConfiguration
}

var config Configuration

// LoadConfigurations load config file
func LoadConfigurations() (err error) {
	config = Configuration{}
	appConfig, err := loadAppConfiguration()
	if err != nil {
		return err
	}
	config.App = appConfig

	loggerConfig, err := loadLoggerConfiguration()
	if err != nil {
		return err
	}
	config.Logger = loggerConfig

	authConfig, err := loadAuthConfiguration()
	if err != nil {
		return err
	}
	config.Auth = authConfig

	return nil
}

func GetConfig() Configuration {
	return config
}
