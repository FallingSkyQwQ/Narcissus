package commands

import (
	"fmt"

	"github.com/FallingSkyQwQ/Narcissus/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check development environment",
	Long:  `Check if all required tools and dependencies are installed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking Narcissus development environment...")
		fmt.Println()

		results := doctor.RunChecks()

		var warnings, errCount int
		for _, result := range results {
			fmt.Printf("%s %s: %s\n", result.Status, result.Name, result.Message)
			if result.Status == doctor.StatusWarning {
				warnings++
				if result.Fix != "" {
					fmt.Printf("  → %s\n", result.Fix)
				}
			} else if result.Status == doctor.StatusError {
				errCount++
				if result.Fix != "" {
					fmt.Printf("  → %s\n", result.Fix)
				}
			}
		}

		fmt.Println()
		if errCount > 0 {
			fmt.Printf("Found %d error(s) and %d warning(s). Please fix the errors before proceeding.\n", errCount, warnings)
		} else if warnings > 0 {
			fmt.Printf("Environment is functional with %d optional improvement(s).\n", warnings)
		} else {
			fmt.Println("All checks passed! Your environment is ready.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
