package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

var postgresDB *sqlx.DB

func getPostgresDB() *sqlx.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s",
		postgresName, postgresName, postgresName)

	var err error
	postgresDB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgresSQL: %v", err)
	}

	postgresDB.SetMaxOpenConns(100)
	postgresDB.SetMaxIdleConns(10)

	return postgresDB
}

// migrationSourcePath is a relative path to the collection storing all migration scripts
const migrationSourcePath = "file://../../../scripts/migrations"
const postgresName = "page_turner_pro"
