package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// DBConnector interface allows dependency injection for testing
type DBConnector interface {
	Open(driverName, dataSourceName string) (*sql.DB, error)
	Ping(db *sql.DB) error
}

// DefaultConnector is the real implementation using sql package
type DefaultConnector struct{}

func (c *DefaultConnector) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

func (c *DefaultConnector) Ping(db *sql.DB) error {
	return db.Ping()
}

// ConnectDB establishes database connection with retry logic and runs migrations
// Maintains backward compatibility with existing code
func ConnectDB() (*sql.DB, error) {
	return ConnectDBWithConnector(&DefaultConnector{})
}

// ConnectDBWithConnector allows dependency injection for testing
func ConnectDBWithConnector(connector DBConnector) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Validate required environment variables
	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("missing required environment variables (DB_HOST, DB_PORT, DB_USER, DB_NAME)")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := connector.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Retry logic with exponential backoff
	maxRetries := 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		lastErr = connector.Ping(db)
		if lastErr == nil {
			break
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(1<<uint(i)) * time.Second
			fmt.Printf("Connection attempt %d failed, retrying in %v...\n", i+1, waitTime)
			time.Sleep(waitTime)
		}
	}

	if lastErr != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, lastErr)
	}

	fmt.Println("Connected to " + dbname)

	// Execute database migrations
	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}
