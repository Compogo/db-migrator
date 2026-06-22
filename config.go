package db_migrator

import (
	"github.com/Compogo/compogo"
)

const (
	// PathFieldName — имя поля для пути к миграциям.
	PathFieldName = "migrator.source"

	// AutoMigrateFieldName — имя поля для автоматического применения миграций.
	AutoMigrateFieldName = "migrator.auto"

	// DriverReplacement — плейсхолдер для подстановки имени драйвера в путь.
	// Используется для автоматического выбора папки с миграциями для конкретной БД.
	DriverReplacement = "{db.driver}"
)

var (
	// PathDefault — путь к миграциям по умолчанию.
	// Пример: file://./migrations/mysql, file://./migrations/postgres
	PathDefault = "file://./migrations/" + DriverReplacement

	// AutoMigrateDefault — автоматическое применение миграций отключено по умолчанию.
	AutoMigrateDefault = false
)

// Config содержит конфигурацию мигратора.
type Config struct {
	Path        string
	AutoMigrate bool
	Driver      string
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Если Driver не задан и зарегистрирован только один драйвер, используется он.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
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
