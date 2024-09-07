package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
	"github.com/surfiniaburger/api-go/cmd/api"
	"github.com/surfiniaburger/api-go/configs"
	"github.com/surfiniaburger/api-go/db"
)

func main() {
	// Load configuration
	cfg := configs.Envs

	// Initialize PostgreSQL connection
	db, err := db.NewPostgresStorage(cfg) // Pass the entire Config struct
	if err != nil {
		log.Fatal(err)
	}

	// Initialize and verify the DB connection
	initStorage(db)

	// Start the API server
	server := api.NewAPIServer(fmt.Sprintf(":%s", cfg.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	// Test the connection with Ping
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected to PostgreSQL!")
}

