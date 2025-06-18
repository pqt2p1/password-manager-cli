package service

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/internal/models"
	"github.com/pqt2p1/password-manager-cli/internal/repository"
	"time"
)

type passwordService struct {
	repo               repository.PasswordRepository
	masterPasswordHash string
}

func NewPasswordService(repo repository.PasswordRepository, masterPasswordHash string) *passwordService {
	return &passwordService{
		repo:               repo,
		masterPasswordHash: masterPasswordHash,
	}
}

func (s *passwordService) SetMasterPassword(masterPassword string) error {
	s.masterPasswordHash = masterPassword
	return nil
}

func (s *passwordService) VerifyMasterPassword(masterPassword string) error {
	if s.masterPasswordHash == "" {
		return fmt.Errorf("master password not set")
	}
	if s.masterPasswordHash != masterPassword {
		return fmt.Errorf("master password does not match")
	}
	return nil
}

func (s *passwordService) AddPassword(site, username, password string) error {
	if s.masterPasswordHash == "" {
		return fmt.Errorf("master password not set")
	}

	store, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("repository failed to load")
	}

	for _, entry := range store.Entries {
		if entry.Site == site && entry.Username == username {
			return fmt.Errorf("password entry already exists")
		}
	}

	newEntry := models.NewPasswordEntry(site, username, password)

	store.Entries = append(store.Entries, *newEntry)

	return s.repo.Save(store)
}

func (s *passwordService) GetPassword(site string) error {
	store, err := s.repo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}

	for i := range store.Entries {
		if store.Entries[i].Site == site {
			return &store.Entries[i], nil
		}
	}

	return nil, fmt.Errorf("no password found for this site: %s", site)
}

func (s *passwordService) ListPassword() (*[]models.PasswordEntry, error) {
	store, err := s.repo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}

	return &store.Entries, nil
}

func (s *passwordService) UpdatePassword(site, username, password string) error {
	store, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	found := false
	for i := range store.Entries {
		if store.Entries[i].Site == site && store.Entries[i].Username == username {
			store.Entries[i].Password = password
			store.Entries[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("password not found for site: %s", site)
	}

	return s.repo.Save(store)
}

func (s *passwordService) DeletePassword(site string) error {
	store, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	indexToDelete := -1
	for i := range store.Entries {
		if store.Entries[i].Site == site {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return fmt.Errorf("password not found for site: %s", site)
	}

	store.Entries = append(
		store.Entries[:indexToDelete],
		store.Entries[indexToDelete+1:]...,
	)
	return s.repo.Save(store)
}
