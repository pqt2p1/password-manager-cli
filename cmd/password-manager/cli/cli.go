package cli

import (
	"fmt"
	"github.com/pqt2p1/password-manager-cli/internal/service"
	"os"
)

type CLI struct {
	service service.PasswordService
}

func NewCLI(service service.PasswordService) *CLI {
	return &CLI{
		service: service,
	}
}

func (c *CLI) printUsage() {
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

func (c *CLI) Run() error {
	if len(os.Args) < 2 {
		c.printUsage()
		return nil
	}

	command := os.Args[1]
	switch command {
	case "add":
		return c.handleAdd()
	case "get":
		return c.handleGet()
	case "list":
		return c.handleList()
	case "delete":
		return c.handleDelete()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		c.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}
