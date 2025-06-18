package main

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/internal/repository"
	"github.com/pqt2p1/password-manager-cli/internal/service"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Setup dependencies
	homeDir, _ := os.UserHomeDir()
	repoPath := filepath.Join(homeDir, ".password-manager", "passwords.json")
	fmt.Println(repoPath)
	repo := repository.NewFileRepository(repoPath)
	svc := service.NewPasswordService(repo)

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(svc)
	case "get":
		handleGet(svc)
	case "list":
		handleList(svc)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}

}

func handleAdd(svc service.PasswordService) {
	if len(os.Args) < 5 {
		fmt.Println("Usage: password-manager add <site> <username> <password>")
	}

	site := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	masterPass := "temp123"
	if err := svc.SetMasterPassword(masterPass); err != nil {
		fmt.Printf("Error setting master password: %v\n", err)
		return
	}

	if err := svc.AddPassword(site, username, password); err != nil {
		fmt.Printf("Error adding password: %v\n", err)
		return
	}

	fmt.Printf("Password added successfully for %s@%s\n", username, site)

}

func handleGet(svc service.PasswordService) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: password-manager get <site>")
	}

	site := os.Args[2]

	entry, err := svc.GetPassword(site)
	if err != nil {
		fmt.Printf("Error getting password: %v\n", err)
		return
	}

	fmt.Println("Site: %s\nUsername: %s\nPassword: %s\n", site, entry.Username, entry.Password)
}

func handleList(svc service.PasswordService) {
	entries, err := svc.ListPassword()
	if err != nil {
		fmt.Printf("Error listing passwords: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No password entries found")
		return
	}

	fmt.Println("Stored passwords:")
	fmt.Println("==================")
	for _, entry := range entries {
		fmt.Printf("Site: %s  | Username: %s  | Created: %s\n", entry.Site, entry.Username, entry.CreatedAt)
	}
}

func printUsage() {
	fmt.Println("Password Manager CLI")
	fmt.Println("====================")
	fmt.Println("Usage:")
	fmt.Println("  password-manager add <site> <username> <password>")
	fmt.Println("  password-manager get <site>")
	fmt.Println("  password-manager list")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  password-manager add github.com john mypassword123")
	fmt.Println("  password-manager get github.com")
	fmt.Println("  password-manager list")
}
