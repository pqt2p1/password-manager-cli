package ui

import "github.com/fatih/color"

var (
	Success = color.New(color.FgGreen).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
)

func SuccessMsg(msg string) string {
	return Success("✅ " + msg)
}

func ErrorMsg(msg string) string {
	return Error("❌" + msg)
}

func InfoMsg(msg string) string {
	return Info("️ℹ️ " + msg)
}

func PasswordPrompt() string {
	return Info("🔐 Enter master password: ")
}
