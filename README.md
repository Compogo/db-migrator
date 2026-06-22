# Compogo DB Migrator

Мигратор для баз данных в фреймворке [Compogo](https://github.com/Compogo/compogo).

На основе [golang-migrate/migrate](https://github.com/golang-migrate/migrate) предоставляет:

* Поддержку различных СУБД (MySQL, PostgreSQL, SQLite)
* Автоматическое применение миграций при старте
* Подстановку драйвера в путь к миграциям
* Плагинную систему драйверов

## Установка

```shell
go get github.com/Compogo/db-migrator
```

## Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/db-migrator"
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithComponents(&db_migrator.Component),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Конфигурация

### Флаги командной строки

```shell
# Путь к миграциям (с поддержкой плейсхолдера {db.driver})
--migrator.source=file://./migrations/{db.driver}

# Автоматическое применение миграций при старте
--migrator.auto=true
```

### Плейсхолдер `{db.driver}`

Путь к миграциям автоматически подставляет имя драйвера:

```shell
# Если указано:
--migrator.source=file://./migrations/{db.driver}

# И драйвер БД = "mysql", то путь будет:
file://./migrations/mysql

# Если драйвер = "postgres", то:
file://./migrations/postgres
```

## Использование

### Ручное управление миграциями

```go
var migrator *migrate.Migrate
container.Invoke(func(m *migrate.Migrate) { migrator = m })

// Применение всех миграций
err := migrator.Up()

// Откат одной миграции
err := migrator.Down()

// Принудительная установка версии
err := migrator.Force(2)
```

### Проверка состояния

```go
version, dirty, err := migrator.Version()
if err != nil {
    log.Fatal(err)
}
log.Printf("Current version: %d, dirty: %v", version, dirty)
```

## Регистрация драйверов

```go
import (
    "github.com/Compogo/db-migrator"
    "github.com/golang-migrate/migrate/v4/database/postgres"
)

func init() {
    db_migrator.Registration("postgres", func(container compogo.Container) (database.Driver, error) {
        var db *sql.DB
        container.Invoke(func(c *sql.DB) { db = c })
        return postgres.WithInstance(db, &postgres.Config{})
    })
}
```

## Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — библиотека миграций
* [Compogo DB Client](https://github.com/Compogo/db-client) — клиент БД

## Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
