package commands

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/internal/scaffold"
	"github.com/spf13/cobra"
)

var initOpts scaffold.Options

var initCmd = &cobra.Command{
	Use:   "init [app-name]",
	Short: "Create a new Narcissus project",
	Long: `Create a new Narcissus project with the specified name.

This command generates a complete project structure including:
- main.go with a sample counter application
- go.mod with required dependencies
- .gitignore configured for Go and Narcissus projects

Example:
  narc init myapp
  narc init myapp --module github.com/user/myapp
  narc init myapp --template basic`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		initOpts.AppName = args[0]

		// 如果没有指定模块名，使用默认值
		if initOpts.ModuleName == "" {
			initOpts.ModuleName = fmt.Sprintf("github.com/example/%s", initOpts.AppName)
		}

		if err := scaffold.CreateProject(initOpts); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initOpts.ModuleName, "module", "", "Go module name (default: github.com/example/[app-name])")
	initCmd.Flags().StringVar(&initOpts.Template, "template", "basic", "Project template to use")
	initCmd.Flags().StringVar(&initOpts.TargetDir, "dir", "", "Target directory (default: [app-name])")
	rootCmd.AddCommand(initCmd)
}
