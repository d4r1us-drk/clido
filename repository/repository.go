package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Repository manages the database connection and migrations for the application.
// It encapsulates the GORM database instance and a migrator responsible for applying database migrations.
type Repository struct {
	db       *gorm.DB  // The GORM database instance
	migrator *Migrator // The migrator responsible for handling database migrations
}

// NewRepository initializes a new Repository instance, setting up the SQLite database connection.
// It also configures a custom GORM logger and applies any pending migrations.
func NewRepository() (*Repository, error) {
	// Determine the database path
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	// Custom logger for GORM, disabling verbose logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Output logger with timestamp
		logger.Config{
			SlowThreshold:             time.Second,   // Log slow SQL queries taking longer than 1 second
			LogLevel:                  logger.Silent, // Disable all log output (silent mode)
			IgnoreRecordNotFoundError: true,          // Ignore record not found errors in logs
			Colorful:                  false,         // Disable colored output in logs
		},
	)

	// Open the SQLite database using GORM
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger, // Use the custom logger
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Initialize the migrator
	migrator := NewMigrator()

	// Create the repository instance
	repo := &Repository{
		db:       db,
		migrator: migrator,
	}

	// Run database migrations
	err = repo.migrator.Migrate(repo.db)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return repo, nil
}

// getDBPath determines the path for the SQLite database based on the operating system.
//
// On Windows, the path is in the APPDATA directory.
// On Unix-based systems, the path is in the ~/.local/share/clido directory.
func getDBPath() (string, error) {
	var dbPath string

	// Determine the correct path based on the operating system
	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			return "", errors.New("the APPDATA environment variable is not set")
		}
		dbPath = filepath.Join(appDataPath, "clido", "data.db")
	} else {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return "", errors.New("the HOME environment variable is not set")
		}
		dbPath = filepath.Join(homePath, ".local", "share", "clido", "data.db")
	}

	// Ensure the database directory exists, creating it if necessary
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("error creating database directory: %w", err)
	}

	return dbPath, nil
}

// Close closes the database connection gracefully.
// It retrieves the underlying SQL database object from GORM and calls its Close method.
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
