package db_migrator

import (
	"errors"

	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
	"github.com/golang-migrate/migrate/v4"
)

// Component — компонент мигратора для Compogo.
// Регистрирует конфигурацию и мигратор в DI-контейнере.
// Автоматически применяет миграции при старте, если AutoMigrate = true.
//
// Пример подключения:
//
//	app.AddComponents(&db_migrator.Component)
//
// Ручной запуск миграций:
//
//	var m *migrate.Migrate
//	container.Invoke(func(migrator *migrate.Migrate) { m = migrator })
//	err := m.Up()
//
// Структура миграций:
//
//	migrations/
//	├── mysql/
//	│   ├── 1_init.up.sql
//	│   └── 1_init.down.sql
//	└── postgres/
//	    ├── 1_init.up.sql
//	    └── 1_init.down.sql
var Component = compogo.Component{
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			NewConfig,
			NewMigrator,
		)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.StringVar(&config.Path, PathFieldName, PathDefault, "path to migrations directory")
			flagSet.BoolVar(&config.AutoMigrate, AutoMigrateFieldName, AutoMigrateDefault, "automatically migrate")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
	Execute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(config *Config, migrator *migrate.Migrate, logger compogo.Logger) error {
			if !config.AutoMigrate {
				return nil
			}

			err := migrator.Up()
			if err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return err
			}

			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}

			version, _, err := migrator.Version()

			logger.GetLogger("Database").GetLogger("migrator").Infof("up to '%d' version", version)

			return nil
		})
	}),
}
