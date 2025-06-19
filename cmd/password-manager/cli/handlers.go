package cli

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/pkg/ui"
	"os"
)

func (c *CLI) handleAdd() error {
	if len(os.Args) < 5 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager add <site> <username> <password>"))

	}

	site := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %s\n", err)))
		return fmt.Errorf(ui.ErrorMsg("Failed to get master password"))
	}

	if err := c.service.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to set master password: %s\n", err)))
		return err
	}

	if err := c.service.AddPassword(site, username, password); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Error adding password: %s\n", err)))
		return err
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Password added successfully for %s@%s\n", username, site)))
	return nil
}

func (c *CLI) handleGet() error {
	if len(os.Args) < 3 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager get <site>"))
		return fmt.Errorf(ui.ErrorMsg("Usage: password-manager get <site>"))
	}

	site := os.Args[2]

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Failed to get master password: %s\n"))
		return err
	}

	if err := c.service.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg("Failed to set master password: %s\n"))
		return err
	}

	entry, err := c.service.GetPassword(site)
	if err != nil {
		fmt.Println(ui.ErrorMsg("Error getting password: %v\n"))
		return err
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Site: %s\nUsername: %s\nPassword: %s\n", site, entry.Username, entry.Password)))
	return nil
}

func (c *CLI) handleList() error {
	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Failed to get master password: %s\n"))
		return err
	}

	if err := c.service.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg("Failed to set master password: %s\n"))
		return err
	}

	entries, err := c.service.ListPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg("Error listing passwords: %v\n"))
		return err
	}

	if len(entries) == 0 {
		fmt.Println(ui.ErrorMsg("No password entries found"))
		return fmt.Errorf(ui.ErrorMsg("No password entries found"))
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
	return nil
}

func (c *CLI) handleDelete() error {
	if len(os.Args) < 3 {
		fmt.Println(ui.ErrorMsg("Usage: password-manager delete <site>"))
		return fmt.Errorf(ui.ErrorMsg("Usage: password-manager delete <site>"))
	}

	site := os.Args[2]

	// Confirm deletion
	fmt.Printf(ui.Warning("⚠️  Delete password for %s? (y/N): "), ui.Bold(site))
	var confirm string
	_, err := fmt.Scanln(&confirm)
	if err != nil {
		fmt.Println(ui.ErrorMsg("Invalid input"))
		return err
	}

	if confirm != "y" && confirm != "Y" {
		fmt.Println(ui.InfoMsg("Deletion cancelled"))
		return err
	}

	masterPass, err := askMasterPassword()
	if err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %v", err)))
		return err
	}

	if err := c.service.SetMasterPassword(masterPass); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to get master password: %v", err)))
		return err
	}

	// Delete the password
	if err := c.service.DeletePassword(site); err != nil {
		fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to delete master password: %v", err)))
		return err
	}

	fmt.Println(ui.SuccessMsg(fmt.Sprintf("Password for %s deleted successfully!", site)))
	return nil
}
