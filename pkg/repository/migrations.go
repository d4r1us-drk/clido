package repository

import (
	"github.com/d4r1us-drk/clido/pkg/models"
	"gorm.io/gorm"
)

type Migration struct {
	ID      uint   `gorm:"primaryKey"`
	Version string `gorm:"uniqueIndex"`
}

type Migrator struct {
	migrations []struct {
		version string
		migrate func(*gorm.DB) error
	}
}

func NewMigrator() *Migrator {
	return &Migrator{
		migrations: []struct {
			version string
			migrate func(*gorm.DB) error
		}{
			{
				version: "1.0",
				migrate: func(db *gorm.DB) error {
					return db.AutoMigrate(&models.Project{}, &models.Task{})
				},
			},

			// Example migration for reference:
			// {
			// 	version: "1.1",
			// 	migrate: func(db *gorm.DB) error {
			// 		return db.Exec("ALTER TABLE projects ADD COLUMN status VARCHAR(50)").Error
			// 	},
			// },
		},
	}
}

func (m *Migrator) Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Migration{})
	if err != nil {
		return err
	}

	var lastMigration Migration
	db.Order("version desc").First(&lastMigration)

	for _, migration := range m.migrations {
		if migration.version > lastMigration.Version {
			err = migration.migrate(db)
			if err != nil {
				return err
			}

			err = db.Create(&Migration{Version: migration.version}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
