package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display detailed version information about the Narcissus CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Narcissus CLI")
		fmt.Printf("  Version:   %s\n", version)
		fmt.Printf("  Go:        %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
