package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FallingSkyQwQ/Narcissus/internal/build"
	"github.com/FallingSkyQwQ/Narcissus/internal/pack"
	"github.com/spf13/cobra"
)

var packageOpts pack.Options

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the application as MSIX",
	Long: `Package the application as MSIX for Windows Store or sideloading.

This command creates an MSIX package from your built executable.
The package can be distributed through the Microsoft Store or sideloaded.

Example:
  narc package
  narc package --version 1.0.1
  narc package --publisher "CN=MyCompany"
  narc package --exe ./dist/myapp.exe`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, _ := cmd.Flags().GetString("dir")
		if projectDir == "" {
			projectDir = "."
		}
		packageOpts.ProjectDir = projectDir

		if packageOpts.AppName == "" {
			packageOpts.AppName = build.GetAppName(projectDir)
		}

		if packageOpts.ExePath == "" {
			packageOpts.ExePath = filepath.Join(packageOpts.OutputDir, packageOpts.AppName+".exe")
		}

		if _, err := os.Stat(packageOpts.ExePath); os.IsNotExist(err) {
			fmt.Printf("Executable not found: %s\n", packageOpts.ExePath)
			fmt.Println("Building first...")
			compress, _ := cmd.Flags().GetBool("compress")
			buildOpts := build.Options{
				ProjectDir: projectDir,
				OutputDir:  packageOpts.OutputDir,
				AppName:    packageOpts.AppName,
				Compress:   compress,
			}
			if err := build.Build(buildOpts); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
		}

		if err := pack.Package(packageOpts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	packageCmd.Flags().String("dir", "", "Project directory (default: current directory)")
	packageCmd.Flags().StringVar(&packageOpts.OutputDir, "output", "dist", "Output directory")
	packageCmd.Flags().StringVar(&packageOpts.AppName, "name", "", "Application name")
	packageCmd.Flags().StringVar(&packageOpts.Version, "version", "1.0.0.0", "Package version")
	packageCmd.Flags().StringVar(&packageOpts.Publisher, "publisher", "CN=Developer", "Publisher CN")
	packageCmd.Flags().StringVar(&packageOpts.ExePath, "exe", "", "Path to executable")
	packageCmd.Flags().Bool("compress", false, "Use UPX compression when building")
	rootCmd.AddCommand(packageCmd)
}
