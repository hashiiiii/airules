package cmd

import (
	"fmt"
	"os"

	"github.com/hashiiiii/airules/pkg/version"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command for airules
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airules",
		Short: "AI Editor rules Installer",
		Long: `airules is a tool for installing rules-for-ai files
for AI-powered editors like Windsurf and Cursor to appropriate locations.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Display help if no flags specified
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying help: %v\n", err)
			}
		},
	}

	// Add subcommands
	cmd.AddCommand(newWindsurfCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

// newVersionCmd returns a command that displays version information
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display version information for airules",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version.Version)
			fmt.Printf("Commit: %s\n", version.Commit)
			fmt.Printf("BuildDate: %s\n", version.BuildDate)
		},
	}
}
