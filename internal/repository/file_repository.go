package repository

import (
	"encoding/json"
	"github.com/pqt2p1/password-manager-cli/internal/models"
	"os"
	"path/filepath"
)

type FileRepository struct {
	filePath string
}

// Constructor
func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{filePath}
}

func (r *FileRepository) Exists() bool {
	_, err := os.Stat(r.filePath)
	return err == nil
}

func (r *FileRepository) Save(store *models.PasswordStore) error {
	// 1. Convert struct -> JSON
	data, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}

	// 2. Create directory if needed
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// 3. Write to file
	return os.WriteFile(r.filePath, data, 0600)
}

func (r *FileRepository) Load() (*models.PasswordStore, error) {
	// 1. Check file exists
	if !r.Exists() {
		return &models.PasswordStore{
			Entries: []models.PasswordEntry{},
		}, nil
	}

	// 2. Read file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var store models.PasswordStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}
