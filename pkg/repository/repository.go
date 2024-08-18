package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/d4r1us-drk/clido/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	repo := &Repository{db: db}
	err = repo.autoMigrate()
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

func (r *Repository) autoMigrate() error {
	return r.db.AutoMigrate(&models.Project{}, &models.Task{})
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
