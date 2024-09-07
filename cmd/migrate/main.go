package main

import (
    "log"
    "os"

    _ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/surfiniaburger/api-go/configs"
    "github.com/surfiniaburger/api-go/configloader"
    "github.com/surfiniaburger/api-go/db"
)

func main() {
    // Load configuration
    configs.Envs = configloader.InitConfig() // Initialize Envs with the config

    // Use configs.Envs as needed
    cfg := configs.Envs

    // Initialize PostgreSQL connection
    db, err := db.NewPostgresStorage(cfg) // Pass the entire config struct
    if err != nil {
        log.Fatal(err)
    }

    // Create a migration driver for PostgreSQL
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Set up migrations with the database instance
    m, err := migrate.NewWithDatabaseInstance(
        "file://cmd/migrate/migrations", // Path to migrations
        "postgres",                      // PostgreSQL database name
        driver,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Check current migration version and dirty state
    v, d, _ := m.Version()
    log.Printf("Version: %d, dirty: %v", v, d)

    // Handle migration commands (up or down)
    cmd := os.Args[len(os.Args)-1]
    if cmd == "up" {
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    }
    if cmd == "down" {
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    }
}
