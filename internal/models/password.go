package models

import (
	"github.com/google/uuid"
	"time"
)

type PasswordEntry struct {
	ID        string    `json:"id"`
	Site      string    `json:"site"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordStore struct {
	MasterPasswordHash string          `json:"master_password_hash"`
	Entries            []PasswordEntry `json:"entries"`
}

func NewPasswordEntry(site, username, password string) *PasswordEntry {
	return &PasswordEntry{
		ID:        uuid.New().String(),
		Site:      site,
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
