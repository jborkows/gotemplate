package db_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for loading migrations from files
	queries "github.com/jborkows/gotemplate/internal/db"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrations", // Path to your migration files
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func TestMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:") // Use in-memory DB for testing
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Add test logic here, e.g., checking if tables were created
}

func TestMigrations2(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	queries := queries.New(db)
	apps, error := queries.ListApps(context.Background())
	// Example test: check if table exists
	if error != nil {
		t.Fatalf("table 'app' does not exist: %v", err)
	}
	if len(apps) == 0 {
		t.Fatal("not found any ")
	}
}
