package cli

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/pkg/ui"
	"golang.org/x/term"
	"syscall"
)

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
