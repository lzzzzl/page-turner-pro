package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func upPostgresDB() (*sqlx.DB, error) {
	db := getPostgresDB()

	// migrate postgres
	localHostPostgresDSN := fmt.Sprintf(
		"postgresql://%s:%s@localhost/%s?sslmode=disable",
		postgresName, postgresName, postgresName,
	)
	m, err := migrate.New(migrationSourcePath, localHostPostgresDSN)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start migration process")
	}

	err = m.Up()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to migrare process")
	}

	return db, nil
}

func downPostgresDB() (*sqlx.DB, error) {
	db := getPostgresDB()

	// migrate postgres
	localHostPostgresDSN := fmt.Sprintf(
		"postgresql://%s:%s@localhost/%s?sslmode=disable",
		postgresName, postgresName, postgresName,
	)
	m, err := migrate.New(migrationSourcePath, localHostPostgresDSN)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start migration process")
	}

	err = m.Down()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to migrare process")
	}

	return db, nil
}

func truncateAllData(t *testing.T, db *sqlx.DB) error {
	template := `
		CREATE OR REPLACE FUNCTION truncate_all_tables() RETURNS void AS $$
		DECLARE
			statements CURSOR FOR
				SELECT tablename FROM %s.pg_catalog.pg_tables
				WHERE schemaname = 'public' AND tablename != 'schema_migrations';
		BEGIN
			FOR stmt IN statements LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' RESTART IDENTITY CASCADE;';
			END LOOP;
		END
		$$ LANGUAGE plpgsql;
		SELECT truncate_all_tables();
	`

	script := fmt.Sprintf(template, postgresName)
	_, err := db.Exec(script)
	if err != nil {
		return errors.WithMessage(err, "failed to truncate all data")
	}

	return nil
}

func initRepository(t *testing.T, db *sqlx.DB, files ...string) (repo *PostgresRepository) {
	// Truncate existing records in all tables
	err := truncateAllData(t, db)
	assert.NoError(t, err)

	// Setup DB again
	loader, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Location(time.UTC),
		// Load predefined Data
		testfixtures.Files(files...),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	require.NoError(t, err)

	err = loader.Load()
	require.NoError(t, err)

	return NewPostgresRepository(context.Background(), db)
}

func TestGetPostgresDB(t *testing.T) {
	db := getPostgresDB()

	assert.NotNil(t, db)
}

func TestUpPostgresDB(t *testing.T) {
	db, err := upPostgresDB()

	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestDownPostgresDB(t *testing.T) {
	db, err := downPostgresDB()

	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestTruncateAllData(t *testing.T) {
	db := getPostgresDB()
	assert.NotNil(t, db)

	err := truncateAllData(t, db)
	assert.NoError(t, err)
}
