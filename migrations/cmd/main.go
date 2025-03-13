package main

import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "log"
)

func main() {
    m, err := migrate.New(
        "file://migrations",
        "postgres://postgres:root@localhost:5433/Linkr?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migrations applied successfully")
}