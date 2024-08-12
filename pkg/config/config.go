package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	DBPath string
}

func NewConfig() *Config {
	return &Config{
		DBPath: getDBPath(),
	}
}

func getDBPath() string {
	var dbPath string

	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			panic("APPDATA environment variable is not set")
		}
		dbPath = filepath.Join(appDataPath, "clido", "data.db")
	} else {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			panic("HOME environment variable is not set")
		}
		dbPath = filepath.Join(homePath, ".local", "share", "clido", "data.db")
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		panic(err)
	}

	return dbPath
}
