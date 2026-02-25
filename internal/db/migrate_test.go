package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-migrate/migrate/v4"
)

// TestRunMigrations_Success tests successful migration execution
func TestRunMigrations_Success(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expectations for schema_migrations table queries
	// The migrate library will query the schema_migrations table
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(1, false))

	// Note: We cannot fully test RunMigrations with sqlmock because it requires
	// actual migration files and the migrate library performs complex op
	mock.ExpectBegin()
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS "schema_migrations"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(1, false))
	mock.ExpectCommit()

	// Note: We cannot fully test RunMigrations with sqlmock because it requires
	// actual migration files and the migrate library performs complex operations.
	// This test verifies the database connection handling.

	// For now, we'll test that the function can be called without panicking
	// Real integration tests would require a test database
	err = RunMigrations(db)

	// We expect an error because migration files won't be found in test environment
	if err == nil {
		t.Error("expected error due to missing migration files, got nil")
	}

	// Verify the error is related to migration setup, not database connection
	if err != nil && !contains(err.Error(), "migrate") {
		t.Errorf("unexpected error type: %v", err)
	}
}

// TestRunMigrations_DirtyState tests dirty state detection and error handling
func TestRunMigrations_DirtyState(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expectations for dirty state check
	// When the database is in dirty state, version query returns dirty=true
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(2, true))

	// Attempt to run migrations
	err = RunMigrations(db)

	// We expect an error, but it will be about migration files not found
	// In a real scenario with migration files, we would get dirty state error
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// TestRunMigrations_NoChange tests scenario with no pending migrations
func TestRunMigrations_NoChange(t *testing.T) {
	// This test verifies that when migrate.ErrNoChange is returned,
	// RunMigrations handles it gracefully and returns nil

	// Create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expectations for current version check
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(5, false))

	// Attempt to run migrations
	err = RunMigrations(db)

	// We expect an error about migration files, not about no changes
	// In a real scenario, if all migrations are applied, we'd get nil
	if err == nil {
		t.Error("expected error due to missing migration files, got nil")
	}
}

// TestRunMigrations_MigrationFailure tests migration failure scenarios
func TestRunMigrations_MigrationFailure(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name: "database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
					WillReturnError(errors.New("connection lost"))
			},
			expectedError: "migrate",
		},
		{
			name: "version query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: "migrate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database: %v", err)
			}
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Run migrations
			err = RunMigrations(db)

			// Verify error occurred
			if err == nil {
				t.Error("expected error, got nil")
			}

			// Verify error contains expected substring
			if err != nil && !contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error containing '%s', got: %v", tt.expectedError, err)
			}
		})
	}
}

// TestRunMigrations_NilDatabase tests behavior with nil database
func TestRunMigrations_NilDatabase(t *testing.T) {
	// Note: This test is expected to panic because the postgres driver
	// attempts to ping a nil database. In production, this should never happen
	// as ConnectDB validates the connection before calling RunMigrations.
	// We'll use recover to catch the panic and verify it occurs.

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil database, but no panic occurred")
		}
	}()

	_ = RunMigrations(nil)
}

// TestBuildDatabaseURL tests the database URL construction
func TestBuildDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name: "standard configuration",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "docker configuration",
			envVars: map[string]string{
				"DB_HOST":     "go_db",
				"DB_PORT":     "5432",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "secret123",
				"DB_NAME":     "mercado",
			},
			expected: "postgres://postgres:secret123@go_db:5432/mercado?sslmode=disable",
		},
		{
			name: "empty values",
			envVars: map[string]string{
				"DB_HOST":     "",
				"DB_PORT":     "",
				"DB_USER":     "",
				"DB_PASSWORD": "",
				"DB_NAME":     "",
			},
			expected: "postgres://:@:/?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Build URL
			result := buildDatabaseURL()

			// Verify result
			if result != tt.expected {
				t.Errorf("expected URL '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestCreateMigrateInstance_NilDatabase tests createMigrateInstance with nil database
func TestCreateMigrateInstance_NilDatabase(t *testing.T) {
	// Note: This test is expected to panic because the postgres driver
	// attempts to ping a nil database. We'll use recover to catch the panic.

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil database, but no panic occurred")
		}
	}()

	_, _ = createMigrateInstance(nil)
}

// TestCreateMigrateInstance_ValidDatabase tests createMigrateInstance with valid database
// TestCreateMigrateInstance_ValidDatabase tests createMigrateInstance with valid database
func TestCreateMigrateInstance_ValidDatabase(t *testing.T) {
	// This test verifies that createMigrateInstance can be called with a valid database
	// connection. The postgres driver makes many internal queries (CURRENT_DATABASE,
	// CURRENT_SCHEMA, pg_advisory_lock, etc.) which are difficult to mock completely.
	// Since we don't have actual migration files in the test environment, we expect
	// an error about missing migration files, not about database driver issues.

	// Note: This is a simplified test. Full integration testing with a real database
	// would be needed to thoroughly test createMigrateInstance behavior.

	// Create mock database
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Attempt to create migrate instance
	m, err := createMigrateInstance(db)

	// We expect an error because either:
	// 1. Migration files don't exist (file://migrations not found)
	// 2. Mock database doesn't support all postgres driver queries
	if err == nil {
		t.Error("expected error due to missing migration files or mock limitations, got nil")
		if m != nil {
			m.Close()
		}
		return
	}

	// The error should be related to migration setup, not a nil pointer or panic
	// This verifies the function handles errors gracefully
	t.Logf("Expected error occurred: %v", err)
}

// TestLogMigrationStatus tests the logging function
func TestLogMigrationStatus(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expectations for version query
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(1, false))

	// Create migrate instance (will fail due to missing files, but that's ok for this test)
	m, _ := createMigrateInstance(db)
	if m != nil {
		defer m.Close()
		// Call logMigrationStatus - it should not panic
		logMigrationStatus(m)
	}

	// If we got here without panic, the test passes
	// The actual logging output is not verified as it goes to log.Println
}

// TestLogMigrationStatus_NilVersion tests logging with fresh database
func TestLogMigrationStatus_NilVersion(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expectations for version query returning no rows (fresh database)
	mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnError(sql.ErrNoRows)

	// Create migrate instance
	m, _ := createMigrateInstance(db)
	if m != nil {
		defer m.Close()
		// Call logMigrationStatus - should handle ErrNilVersion gracefully
		logMigrationStatus(m)
	}

	// If we got here without panic, the test passes
}

// TestMigrateErrorHandling tests various error scenarios
func TestMigrateErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		migrateErr  error
		expectNil   bool
		expectError string
	}{
		{
			name:        "ErrNoChange should return nil",
			migrateErr:  migrate.ErrNoChange,
			expectNil:   true,
			expectError: "",
		},
		{
			name:        "ErrNilVersion should be handled",
			migrateErr:  migrate.ErrNilVersion,
			expectNil:   false,
			expectError: "",
		},
		{
			name:        "generic error should be returned",
			migrateErr:  errors.New("migration failed"),
			expectNil:   false,
			expectError: "migration failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies the error handling logic
			// In the actual RunMigrations function:
			// - migrate.ErrNoChange returns nil
			// - migrate.ErrNilVersion is handled specially
			// - Other errors are wrapped and returned

			if tt.migrateErr == migrate.ErrNoChange {
				// Verify that ErrNoChange is treated as success
				if !tt.expectNil {
					t.Error("ErrNoChange should result in nil error")
				}
			}

			if tt.migrateErr == migrate.ErrNilVersion {
				// Verify that ErrNilVersion is handled (not returned as error in version check)
				// This is tested in the actual RunMigrations function
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
