package main

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/internal/repository"
	"github.com/pqt2p1/password-manager-cli/internal/service"
	"github.com/pqt2p1/password-manager-cli/pkg/ui"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"syscall"
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
	case "delete":
		handleDelete(svc)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}

}

func askMasterPassword() (string, error) {
	fmt.Print(ui.PasswordPrompt())
	// Hide password input
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println()
	return string(bytePassword), nil

}

func handleAdd(svc service.PasswordService) {
	if len(os.Args) < 5 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager add <site> <username> <password>"))
		return
	}

	site := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %s\n", err)))
		return
	}

	if err := svc.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to set master password: %s\n", err)))
		return
	}

	if err := svc.AddPassword(site, username, password); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Error adding password: %s\n", err)))
		return
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Password added successfully for %s@%s\n", username, site)))
}

func handleGet(svc service.PasswordService) {
	if len(os.Args) < 3 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager get <site>"))
		return
	}

	site := os.Args[2]

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Failed to get master password: %s\n"))
		return
	}

	if err := svc.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg("Failed to set master password: %s\n"))
		return
	}

	entry, err := svc.GetPassword(site)
	if err != nil {
		fmt.Println(ui.ErrorMsg("Error getting password: %v\n"))
		return
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Site: %s\nUsername: %s\nPassword: %s\n", site, entry.Username, entry.Password)))

}

func handleList(svc service.PasswordService) {
	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Failed to get master password: %s\n"))
		return
	}

	if err := svc.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg("Failed to set master password: %s\n"))
		return
	}

	entries, err := svc.ListPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Error listing passwords: %v\n"))
		return
	}

	if len(entries) == 0 {
		fmt.Println(ui.ErrorMsg("No password entries found"))
		fmt.Println("No password entries found")
		return
	}

	fmt.Println(ui.Bold("\n📋 Stored Passwords"))
	fmt.Println(ui.Bold("=================="))

	for i, entry := range entries {
		fmt.Printf("%s %s | %s | %s | %s\n",
			ui.Info(fmt.Sprintf("%d.", i+1)),
			ui.Bold(entry.Site),
			ui.Success(entry.Username),
			ui.Warning(entry.Password),
			ui.Info(entry.CreatedAt.Format("2006-01-02")),
		)
	}

	fmt.Printf("\n%s\n", ui.InfoMsg(fmt.Sprintf("Total: %d entries", len(entries))))
}

func handleDelete(svc service.PasswordService) {
	if len(os.Args) < 3 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager delete <site>"))
		return
	}

	site := os.Args[2]

	// Confirm deletion
	fmt.Printf(ui.Warning("⚠️  Delete password for %s? (y/N): "), ui.Bold(site))
	var confirm string
	_, err := fmt.Scanln(&confirm)
	if err != nil {
		fmt.Println(ui.ErrorMsg("Invalid input"))
		return
	}

	if confirm != "y" && confirm != "Y" {
		fmt.Println(ui.InfoMsg("Deletion cancelled"))
		return
	}

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %v", err)))
		return
	}

	if err := svc.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %v", err)))
		return
	}

	// Delete the password
	if err := svc.DeletePassword(site); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to delete master password: %v", err)))
		return
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Password for %s deleted successfully!", site)))
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
	fmt.Println("  password-manager add github.com john")
	fmt.Println("  password-manager get github.com")
	fmt.Println("  password-manager list")
}
