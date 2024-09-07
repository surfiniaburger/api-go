package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
	"github.com/surfiniaburger/api-go/configloader" // Use configloader instead of configs
)

// NewPostgresStorage initializes a new database connection to PostgreSQL
func NewPostgresStorage(cfg configloader.Config) (*sql.DB, error) {
	// Use the DBURL from the config struct
	db, err := sql.Open("pgx", cfg.DBURL)
	if err != nil {
		return nil, err
	}

	// Optionally, you can add a ping to check the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database!")
	return db, nil
}
