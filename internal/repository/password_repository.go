package repository

import "github.com/pqt2p1/password-manager-cli/internal/models"

type PasswordRepository interface {
	// Save entire password store to storage
	Save(store *models.PasswordEntry)
	Load(store *models.PasswordEntry)
	Exists() bool
}
