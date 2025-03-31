package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/version"
	"github.com/spf13/cobra"
)

// newVersionCmd returns a command that displays version information.
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
