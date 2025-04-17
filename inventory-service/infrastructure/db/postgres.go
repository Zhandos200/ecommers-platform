package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgres() *sqlx.DB {
	// Optional: load .env only if running outside Docker
	_ = godotenv.Load() // ignore error silently, fallback to os.Getenv

	// Build the connection string from environment variables
	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	log.Println("Connecting with DSN:", dsn) // Temporary log for debugging

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("DB Connection error:", err)
	}
	log.Println("Successfully connected to PostgreSQL")
	return db
}
