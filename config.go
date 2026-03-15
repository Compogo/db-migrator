package db_migrator

import "github.com/Compogo/compogo/configurator"

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
	Driver      Driver
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

	return config
}
