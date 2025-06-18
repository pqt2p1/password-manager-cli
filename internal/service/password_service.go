package service

import "github.com/pqt2p1/password-manager-cli/internal/models"

type PasswordService interface {
	AddPassword(site, username, password string) error
	GetPassword(site string) (*models.PasswordStore, error)
	ListPassword() (*[]models.PasswordStore, error)
	UpdatePassword(site, username, password string) error
	DeletePassword(site string) error

	SetMasterPassword(masterPassword string) error
	VerifyMasterPassword(masterPassword string) error
}
