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
    "gorm.io/gorm/logger"
	"gorm.io/gorm"
)

type Repository struct {
	db       *gorm.DB
	migrator *Migrator
}

func NewRepository() (*Repository, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

    // Custom logger for GORM, we use this to disable GORM's verbose messages
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second,   // Slow SQL threshold
            LogLevel:                  logger.Silent, // Log level
            IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
            Colorful:                  false,         // Disable color
        },
    )

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: newLogger,
    })
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	migrator := NewMigrator()

	repo := &Repository{
		db:       db,
		migrator: migrator,
	}

	err = repo.migrator.Migrate(repo.db)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return repo, nil
}

func getDBPath() (string, error) {
	var dbPath string

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

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("error creating database directory: %w", err)
	}

	return dbPath, nil
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
