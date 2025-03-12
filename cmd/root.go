package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command for airules
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airules",
		Short: "AI Editor rules Installer",
		Long:  "airules is a tool for installing rules-for-ai files for AI-powered editors like Windsurf and Cursor to appropriate locations.",
		Run: func(cmd *cobra.Command, args []string) {
			// Display help if no flags specified
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying help: %v\n", err)
			}
		},

		// Add custom error handling
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Disable completion command
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands
	cmd.AddCommand(newWindsurfCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newRulesCmd())
	cmd.AddCommand(newInitCmd())

	return cmd
}
