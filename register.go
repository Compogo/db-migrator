package db_migrator

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/types/linker"
	"github.com/golang-migrate/migrate/v4/database"
)

// getters — хранилище зарегистрированных драйверов мигратора.
// Ключ — имя драйвера, значение — функция создания database.Driver.
var getters = linker.NewLinker[string, Getter](linker.KeyStringNormalizer[Getter]())

// Registration регистрирует драйвер мигратора.
// Должна вызываться в init() каждого пакета драйвера.
//
// Пример регистрации MySQL-драйвера:
//
//	func init() {
//	    db_migrator.Registration("mysql", func(container compogo.Container) (database.Driver, error) {
//	        var db *sql.DB
//	        container.Invoke(func(c *sql.DB) { db = c })
//	        return mysql.New(db), nil
//	    })
//	}
func Registration(driverName string, getter Getter) {
	getters.Add(driverName, getter)
}

// Getter — фабричная функция для создания database.Driver.
// Принимает DI-контейнер для получения зависимостей драйвера.
type Getter func(container compogo.Container) (database.Driver, error)
