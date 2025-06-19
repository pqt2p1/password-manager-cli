package main

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/cmd/password-manager/cli"
	"github.com/pqt2p1/password-manager-cli/internal/repository"
	"github.com/pqt2p1/password-manager-cli/internal/service"
	"os"
	"path/filepath"
)

func main() {
	// Setup dependencies
	homeDir, _ := os.UserHomeDir()
	repoPath := filepath.Join(homeDir, ".password-manager", "passwords.json")

	repo := repository.NewFileRepository(repoPath)
	svc := service.NewPasswordService(repo)

	// Create CLI
	app := cli.NewCLI(svc)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
