package repository

import (
	"github.com/d4r1us-drk/clido/models"
	"gorm.io/gorm"
)

// Migration represents a database migration entry.
// Each migration is uniquely identified by a version string.
type Migration struct {
	ID      uint   `gorm:"primaryKey"`  // Primary key for the migration
	Version string `gorm:"uniqueIndex"` // Unique version identifier for each migration
}

// Migrator is responsible for applying database migrations.
// It holds a list of migrations, each associated with a version and a function that applies the migration.
type Migrator struct {
	migrations []struct {
		version string               // The version of the migration
		migrate func(*gorm.DB) error // The function that applies the migration
	}
}

// NewMigrator initializes a new Migrator with a list of migrations.
//
// Each migration is represented by a version string and a function that performs the migration.
// In this case, the initial migration (version "1.0") creates the `Project` and `Task` tables.
func NewMigrator() *Migrator {
	return &Migrator{
		migrations: []struct {
			version string
			migrate func(*gorm.DB) error
		}{
			{
				version: "1.0", // The first version of the database schema
				migrate: func(db *gorm.DB) error {
					// Automatically migrates the schema for the Project and Task models
					return db.AutoMigrate(&models.Project{}, &models.Task{})
				},
			},
			// Example of how to add a new migration:
			// {
			//   version: "1.1",
			//   migrate: func(db *gorm.DB) error {
			//     // SQL or schema changes for version 1.1
			//     return db.Exec("ALTER TABLE projects ADD COLUMN status VARCHAR(50)").Error
			//   },
			// },
		},
	}
}

// Migrate applies any pending migrations to the database.
//
// It first ensures that the `Migration` table exists, then checks the latest applied migration.
// Migrations that have a version greater than the last applied one are executed sequentially.
// After each migration is applied, a record is inserted into the `Migration` table.
func (m *Migrator) Migrate(db *gorm.DB) error {
	// Ensure the Migration table exists
	err := db.AutoMigrate(&Migration{})
	if err != nil {
		return err
	}

	// Retrieve the latest migration version from the database
	var lastMigration Migration
	db.Order("version desc").First(&lastMigration)

	// Apply pending migrations
	for _, migration := range m.migrations {
		if migration.version > lastMigration.Version {
			// Execute the migration function
			err = migration.migrate(db)
			if err != nil {
				return err
			}

			// Record the applied migration version
			err = db.Create(&Migration{Version: migration.version}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
