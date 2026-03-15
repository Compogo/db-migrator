# Compogo DB Migrator 📦

**Compogo DB Migrator** — компонент для управления миграциями базы данных, построенный на базе [golang-migrate](https://github.com/golang-migrate/migrate). Полностью интегрируется с экосистемой Compogo.

## 🚀 Установка

```shell
go get github.com/Compogo/db-migrator
```


### 📦 Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/db-client"
    "github.com/Compogo/db-migrator"
    _ "github.com/Compogo/postgres" // ваш драйвер БД
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithOsSignalCloser(),
        db_client.Component,      // выбираем драйвер через --db.driver
        db_migrator.Component,    // мигратор автоматически подхватит тот же драйвер
        compogo.WithComponents(
            // ... ваши компоненты
        ),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

### ✨ Возможности

#### 🔌 Плагинная архитектура драйверов

Драйверы БД сами регистрируют свою реализацию для мигратора:

```go
// В драйвере postgres
func init() {
    db_migrator.Registration(Postgres, NewPostgresMigrationDriver)
}
```

#### 📁 Умный путь к миграциям

Путь к файлам миграций может содержать плейсхолдер `{db.driver}`:

```shell
--migrator.source="file://./migrations/{db.driver}"
```

При использовании драйвера `postgres` путь автоматически станет `file://./migrations/postgres.`

#### 🚦 Автоматические миграции

```shell
./myapp --migrator.auto=true
```

Миграции выполняются в фазе `Execute` до запуска основных сервисов.

#### ⏱️ Бесконечный таймаут

Миграции — критическая операция, которую нельзя прерывать. Компонент устанавливает таймаут = 0, позволяя миграциям выполняться столько, сколько нужно.

### 🔗 Интеграция с драйверами

Драйверы должны зарегистрировать свою реализацию `database.Driver` для мигратора:

```go
type PostgresMigrationDriver struct {
    // ...
}

func NewPostgresMigrationDriver(container container.Container) (database.Driver, error) {
    var db *sql.DB
    if err := container.Invoke(func(sqlDB *sql.DB) { db = sqlDB }); err != nil {
        return nil, err
    }
    return postgres.WithConnection(db), nil
}

func init() {
    db_migrator.Registration(Postgres, NewPostgresMigrationDriver)
}
```
