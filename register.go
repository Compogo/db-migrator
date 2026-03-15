package db_migrator

import (
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/db-client/driver"
	"github.com/Compogo/types/linker"
	"github.com/golang-migrate/migrate/v4/database"
)

var (
	// getters stores constructor functions for each registered driver.
	// The linker associates each Driver with its corresponding Getter.
	getters = linker.NewLinker[driver.Driver, Getter]()
)

// Registration registers a new database driver constructor for migrations.
// This function should be called during driver package initialization.
// The driver will then be available for use via the migrator component.
//
// Example (in postgres driver):
//
//	func init() {
//	    db_migrator.Registration(Postgres, NewPostgresMigrationDriver)
//	}
func Registration(d driver.Driver, getter Getter) {
	getters.Add(d, getter)
}

// Getter is a function type that creates a new migration driver instance.
// It receives the DI container which may contain dependencies like config or logger,
// and returns a database.Driver compatible with golang-migrate.
type Getter func(container container.Container) (database.Driver, error)
