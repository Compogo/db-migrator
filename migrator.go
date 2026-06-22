package db_migrator

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// NewMigrator создаёт новый экземпляр мигратора.
// Поддерживает подстановку имени драйвера в путь через плейсхолдер {db.driver}.
//
// Примеры путей:
//   - file://./migrations/{db.driver} → file://./migrations/mysql
//   - file://./migrations/mysql → file://./migrations/mysql (без подстановки)
func NewMigrator(config *Config, container compogo.Container, logger compogo.Logger) (*migrate.Migrate, error) {
	path := strings.ReplaceAll(config.Path, DriverReplacement, config.Driver)

	logger = logger.GetLogger("Database").GetLogger("migrator")
	logger.Infof("migration path '%s'", path)

	getter, err := getters.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[Database][migrator] driver '%s' getter undefined: %w", config.Driver, err)
	}

	driver, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[Database][migrator] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.Infof("[db-migrator] usage sql driver '%s'", config.Driver)

	return migrate.NewWithDatabaseInstance(path, strings.ToLower(config.Driver), driver)
}
