package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// buildDatabaseURL constructs PostgreSQL connection string from environment variables
func buildDatabaseURL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
}

// createMigrateInstance creates and configures a migrate instance
func createMigrateInstance(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// logMigrationStatus logs the current migration state
func logMigrationStatus(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Println("[MIGRATE] Current version: none (fresh database)")
		} else {
			log.Printf("[MIGRATE] Failed to get version: %v", err)
		}
		return
	}

	log.Printf("[MIGRATE] Current version: %d, Dirty: %t", version, dirty)
}

// RunMigrations executes all pending database migrations
// Returns error if migrations fail or database is in dirty state
func RunMigrations(db *sql.DB) error {
	log.Println("[MIGRATE] Starting database migrations")

	// Create migrate instance
	m, err := createMigrateInstance(db)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Log current migration status
	logMigrationStatus(m)

	// Check for dirty state before executing migrations
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state at version %d, manual intervention required", version)
	}

	// Execute pending migrations
	log.Println("[MIGRATE] Executing pending migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("[MIGRATE] No pending migrations found")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	// Log completion with final version
	finalVersion, _, err := m.Version()
	if err != nil {
		log.Println("[MIGRATE] Database migrations completed successfully")
	} else {
		log.Printf("[MIGRATE] Database migrations completed successfully (version: %d)", finalVersion)
	}

	return nil
}

// RollbackMigrations rolls back the specified number of migrations
// Used for development and testing purposes
func RollbackMigrations(db *sql.DB, steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be a positive integer, got: %d", steps)
	}

	log.Printf("[MIGRATE] Starting rollback of %d migration(s)", steps)

	// Create migrate instance
	m, err := createMigrateInstance(db)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Log current migration status
	logMigrationStatus(m)

	// Get current version before rollback
	versionBefore, _, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return fmt.Errorf("no migrations to rollback (database is at version 0)")
		}
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Execute rollback using Steps (negative value for down migrations)
	log.Printf("[MIGRATE] Executing rollback from version %d...", versionBefore)
	if err := m.Steps(-steps); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("[MIGRATE] No migrations to rollback")
			return nil
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Log completion with new version
	versionAfter, _, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Println("[MIGRATE] Rollback completed successfully (database at version 0)")
		} else {
			log.Println("[MIGRATE] Rollback completed successfully")
		}
	} else {
		log.Printf("[MIGRATE] Rollback completed successfully (version: %d -> %d)", versionBefore, versionAfter)
	}

	return nil
}

// GetCurrentVersion returns the current migration version and dirty state
// Returns (version, dirty, error)
// If database has no migrations applied, returns (0, false, nil)
func GetCurrentVersion(db *sql.DB) (uint, bool, error) {
	// Create migrate instance
	m, err := createMigrateInstance(db)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version and dirty state
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// No migrations applied yet, return version 0
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get version: %w", err)
	}

	return version, dirty, nil
}
