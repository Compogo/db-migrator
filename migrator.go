package db_migrator

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// NewMigrator creates a new golang-migrate instance based on the configuration.
// It performs the following steps:
//  1. Replaces {db.driver} placeholder in the path with actual driver name
//  2. Looks up the Getter for the configured driver
//  3. Creates the migration driver instance
//  4. Returns a configured migrate.Migrate instance
//
// The migrator is not started automatically — use the Component for lifecycle management.
func NewMigrator(config *Config, container container.Container, informer logger.Informer) (*migrate.Migrate, error) {
	path := strings.ReplaceAll(config.Path, DriverReplacement, config.Driver.String())

	informer.Infof("[db-migrator] migration path '%s'", path)

	getter, err := getters.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[db-migrator] driver '%s' getter undefined: %w", config.Driver, err)
	}

	driver, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[db-migrator] driver '%s' create failed: %w", config.Driver, err)
	}

	informer.Infof("[db-migrator] usage sql driver '%s'", config.Driver.String())

	return migrate.NewWithDatabaseInstance(path, strings.ToLower(config.Driver.String()), driver)
}
