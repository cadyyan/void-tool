package configuration

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Logging LoggingConfiguration
	HTTP    HTTPConfiguration
	RS      RunescapeConfiguration
	SQLite  SQLiteConfiguration
}

func NewConfigurationFromEnv() (Configuration, error) {
	var config Configuration
	if err := envconfig.Process("void", &config); err != nil {
		return Configuration{}, fmt.Errorf(
			"unable to load configuration from environment: %w",
			err,
		)
	}

	return config, nil
}
