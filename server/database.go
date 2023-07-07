package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/TurniXXD/go-rss/internal/database"

	_ "github.com/lib/pq"
)

func initDb() *database.Queries {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/go-rss-db?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	queries := database.New(conn)
	apiCfg = apiConfig{DB: queries}
	return queries
}
