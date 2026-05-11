package db_migrator

import (
	"errors"

	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/compogo/logger"
	"github.com/golang-migrate/migrate/v4"
)

// Component is a ready-to-use Compogo component that provides database migrations.
// It automatically:
//   - Registers Config and Migrator in the DI container
//   - Adds command-line flags for migration configuration
//   - Applies configuration during Configuration phase
//   - Runs migrations (if enabled) during Execute phase with no timeout
//
// Usage:
//
//	compogo.WithComponents(
//	    db_client.Component,      // database client (provides driver name)
//	    db_migrator.Component,    // migrations
//	    // ... driver components (postgres, mysql, etc.)
//	)
//
// The driver name is automatically taken from db-client configuration
// and used to select the appropriate migration driver.
var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewMigrator,
		)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.StringVar(&config.Path, PathFieldName, PathDefault, "path to migrations directory")
			flagSet.BoolVar(&config.AutoMigrate, AutoMigrateFieldName, AutoMigrateDefault, "automatically migrate")
		})
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
	Execute: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(config *Config, migrator *migrate.Migrate, informer logger.Informer) error {
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

			informer.Infof("[db-migrator]: up to '%d' version", version)

			return nil
		})
	}),
}
