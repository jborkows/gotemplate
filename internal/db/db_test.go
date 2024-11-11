package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for loading migrations from files
	queries "github.com/jborkows/gotemplate/internal/db"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
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

func TestConnectionString(t *testing.T) {
	tempFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON&_cache_size=2000", tempFile.Name())
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	// Step 3: Test the database by creating a table
	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		return
	}

	fmt.Println("Table 'test' created successfully.")

	// Step 4: Verify the journal mode and foreign keys settings
	var journalMode, foreignKeys string
	var cacheSize int
	db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	db.QueryRow("PRAGMA cache_size").Scan(&cacheSize)

	fmt.Printf("Journal mode: %s\n", journalMode)
	fmt.Printf("Foreign keys enabled: %s\n", foreignKeys)
	fmt.Printf("Cache size: %d pages\n", cacheSize)
	assert.Equal(t, "wal", journalMode)
	assert.Equal(t, "1", foreignKeys)
	assert.Equal(t, 2000, cacheSize)
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
