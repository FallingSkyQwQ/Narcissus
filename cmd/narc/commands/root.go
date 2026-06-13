package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "0.4.0"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "narc",
	Short: "Narcissus CLI - Build WinUI apps with Go",
	Long: fmt.Sprintf(`Narcissus is a Go framework for building Windows UI applications.

Version: %s

This CLI provides tools for project scaffolding, development, building, and packaging.

Quick Start:
  narc init myapp    # Create a new project
  cd myapp
  narc run           # Run with hot reload
  narc build         # Build for production
  narc package       # Create MSIX package

Commands:
  init      Create a new Narcissus project
  run       Run with hot reload (requires air)
  build     Build production executable
  package   Create MSIX package
  doctor    Check development environment
  version   Show version information

For more information, visit: https://github.com/FallingSkyQwQ/Narcissus`, version),
	Version: version,
}

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 子命令通过各自的 init() 函数注册
}
