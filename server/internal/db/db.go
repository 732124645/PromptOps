package db

import (
	"os"
	"path/filepath"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Open opens (creating if needed) the SQLite database and runs migrations.
func Open(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	d, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := d.AutoMigrate(
		&models.Prompt{}, &models.PromptVersion{}, &models.Agent{},
		&models.Workflow{}, &models.AuditLog{}, &models.RunLog{},
	); err != nil {
		return nil, err
	}
	return d, nil
}
