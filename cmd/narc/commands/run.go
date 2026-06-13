package commands

import (
	"fmt"
	"os"

	"github.com/FallingSkyQwQ/Narcissus/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application in development mode with hot reload",
	Long: `Run the application in development mode with hot reload.

This command uses 'air' to watch for file changes and automatically
rebuild and restart your application.

The first run will create a .air.toml configuration file if it doesn't exist.

Example:
  narc run
  narc run --dir ./myapp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, _ := cmd.Flags().GetString("dir")
		if projectDir == "" {
			projectDir = "."
		}
		mainFile := projectDir + "/main.go"
		if _, err := os.Stat(mainFile); os.IsNotExist(err) {
			return fmt.Errorf("not a Narcissus project: %s/main.go not found", projectDir)
		}
		if err := runner.RunWithAir(projectDir); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	runCmd.Flags().String("dir", "", "Project directory (default: current directory)")
	rootCmd.AddCommand(runCmd)
}
