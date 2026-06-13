package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "narc",
	Short: "Narcissus CLI - Build WinUI apps with Go",
	Long: `Narcissus is a Go framework for building Windows UI applications.

This CLI provides tools for project scaffolding, development, building, and packaging.`,
	Version: "0.4.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 子命令将在后续步骤中添加
}
