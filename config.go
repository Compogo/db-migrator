package db_migrator

import (
	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/db-client/driver"
)

const (
	PathFieldName        = "migrator.source"
	AutoMigrateFieldName = "migrator.auto"

	DriverReplacement = "{db.driver}"

	PathDefault        = "file://./migrations/" + DriverReplacement
	AutoMigrateDefault = false
)

type Config struct {
	Path        string
	AutoMigrate bool
	Driver      driver.Driver
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) *Config {
	if config.Path == "" || config.Path == PathDefault {
		configurator.SetDefault(PathFieldName, PathDefault)
		config.Path = configurator.GetString(PathFieldName)
	}

	if config.AutoMigrate || config.AutoMigrate == AutoMigrateDefault {
		configurator.SetDefault(AutoMigrateFieldName, AutoMigrateDefault)
		config.AutoMigrate = configurator.GetBool(AutoMigrateFieldName)
	}

	if config.Driver == "" && getters.Len() == 1 {
		config.Driver = getters.Keys()[0]
	}

	return config
}
