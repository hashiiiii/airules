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
		Long:  "airules is a tool for installing rules-for-ai files for AI-powered editors to appropriate locations.",
		Run: func(cmd *cobra.Command, args []string) {
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
	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newTemplateCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInitCmd())

	return cmd
}
