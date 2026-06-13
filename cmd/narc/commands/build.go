package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FallingSkyQwQ/Narcissus/internal/build"
	"github.com/spf13/cobra"
)

var buildOpts build.Options

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application for production",
	Long: `Build the application for production.

This command compiles your Narcissus application into a single executable file.
By default, it uses UPX compression to reduce file size.

Example:
  narc build
  narc build --output ./dist
  narc build --name myapp --no-compress
  narc build --ldflags "-s -w -X main.version=1.0.0"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, _ := cmd.Flags().GetString("dir")
		if projectDir == "" {
			projectDir = "."
		}
		buildOpts.ProjectDir = projectDir

		mainFile := filepath.Join(projectDir, "main.go")
		if _, err := os.Stat(mainFile); os.IsNotExist(err) {
			return fmt.Errorf("not a Narcissus project: %s/main.go not found", projectDir)
		}

		if buildOpts.AppName == "" {
			buildOpts.AppName = build.GetAppName(projectDir)
		}

		noCompress, _ := cmd.Flags().GetBool("no-compress")
		if noCompress {
			buildOpts.Compress = false
		}

		if err := build.Build(buildOpts); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	buildCmd.Flags().String("dir", "", "Project directory (default: current directory)")
	buildCmd.Flags().StringVar(&buildOpts.OutputDir, "output", "dist", "Output directory")
	buildCmd.Flags().StringVar(&buildOpts.AppName, "name", "", "Application name (default: from go.mod)")
	buildCmd.Flags().BoolVar(&buildOpts.Compress, "compress", true, "Use UPX compression")
	buildCmd.Flags().StringVar(&buildOpts.LDFlags, "ldflags", "", "Additional linker flags")
	buildCmd.Flags().Bool("no-compress", false, "Disable UPX compression")
	rootCmd.AddCommand(buildCmd)
}
