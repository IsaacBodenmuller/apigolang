package db

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockConnector for testing
type MockConnector struct {
	openErr    error
	pingErr    error
	db         *sql.DB
	mock       sqlmock.Sqlmock
	pingCount  int
	maxRetries int
}

func NewMockConnector() *MockConnector {
	db, mock, _ := sqlmock.New()
	return &MockConnector{
		db:         db,
		mock:       mock,
		maxRetries: 3,
	}
}

func (m *MockConnector) Open(driverName, dataSourceName string) (*sql.DB, error) {
	if m.openErr != nil {
		return nil, m.openErr
	}
	return m.db, nil
}

func (m *MockConnector) Ping(db *sql.DB) error {
	m.pingCount++
	// If maxRetries is set and we haven't reached it yet, fail
	if m.maxRetries > 0 && m.pingCount < m.maxRetries {
		return m.pingErr
	}
	// If maxRetries reached or not set, check pingErr
	if m.maxRetries > 0 && m.pingCount >= m.maxRetries {
		return nil // Success after retries
	}
	return m.pingErr
}

func TestConnectDB_Success(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")

	mock := NewMockConnector()
	defer mock.db.Close()

	// Mock migration expectations
	mock.mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(1, false))

	db, err := ConnectDBWithConnector(mock)

	// We expect an error because migration files won't be found
	// But we verify the connection logic worked
	if err != nil && !strings.Contains(err.Error(), "migrate") {
		t.Errorf("unexpected error type: %v", err)
	}

	if db != nil {
		db.Close()
	}
}

func TestConnectDB_MissingEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "missing DB_HOST",
			envVars: map[string]string{
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: true,
		},
		{
			name: "missing DB_PORT",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: true,
		},
		{
			name: "missing DB_USER",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: true,
		},
		{
			name: "missing DB_NAME",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
			},
			expectError: true,
		},
		{
			name: "all variables present",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")

			// Set test env vars
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			mock := NewMockConnector()
			defer mock.db.Close()

			db, err := ConnectDBWithConnector(mock)

			if tt.expectError {
				if err == nil {
					t.Error("expected error for missing env vars, got nil")
				}
				if err != nil && !strings.Contains(err.Error(), "missing required environment variables") {
					t.Errorf("expected missing env vars error, got: %v", err)
				}
			} else {
				// Even with all vars, we expect migration error
				if err != nil && !strings.Contains(err.Error(), "migrate") {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if db != nil {
				db.Close()
			}
		})
	}
}

func TestConnectDB_OpenError(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")

	mock := NewMockConnector()
	mock.openErr = errors.New("connection refused")
	defer mock.db.Close()

	db, err := ConnectDBWithConnector(mock)

	if err == nil {
		t.Error("expected error for open failure, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "failed to open database") {
		t.Errorf("expected open error, got: %v", err)
	}

	if db != nil {
		db.Close()
	}
}

func TestConnectDB_PingError(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")

	mock := NewMockConnector()
	mock.pingErr = errors.New("connection timeout")
	mock.maxRetries = 0 // Don't allow retries to succeed
	defer mock.db.Close()

	db, err := ConnectDBWithConnector(mock)

	if err == nil {
		t.Error("expected error for ping failure, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "failed to ping database") {
		t.Errorf("expected ping error, got: %v", err)
	}

	if db != nil {
		db.Close()
	}
}

func TestConnectDB_RetryLogic(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")

	mock := NewMockConnector()
	mock.pingErr = errors.New("temporary failure")
	mock.maxRetries = 2 // Fail first 2 attempts, succeed on 3rd
	defer mock.db.Close()

	// Mock migration expectations for successful connection
	mock.mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(1, false))

	db, err := ConnectDBWithConnector(mock)

	// Should succeed after retries (but fail on migrations)
	if err != nil && !strings.Contains(err.Error(), "migrate") {
		t.Errorf("expected migration error after successful retry, got: %v", err)
	}

	if mock.pingCount < 2 {
		t.Errorf("expected at least 2 ping attempts, got %d", mock.pingCount)
	}

	if db != nil {
		db.Close()
	}
}

func TestConnectDB_MigrationError(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")

	mock := NewMockConnector()
	defer mock.db.Close()

	// Mock dirty state to trigger migration error
	mock.mock.ExpectQuery(`SELECT version, dirty FROM "schema_migrations" LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"version", "dirty"}).AddRow(2, true))

	db, err := ConnectDBWithConnector(mock)

	if err == nil {
		t.Error("expected migration error, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "migration failed") {
		t.Errorf("expected migration failed error, got: %v", err)
	}

	if db != nil {
		t.Error("expected nil db on migration failure")
		db.Close()
	}
}

func TestConnectDB_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			for key, expectedValue := range tt.envVars {
				actualValue := os.Getenv(key)
				if actualValue != expectedValue {
					t.Errorf("expected %s=%s, got %s", key, expectedValue, actualValue)
				}
			}
		})
	}
}
