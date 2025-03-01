package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

// newWindsurfCmd returns the windsurf command
func newWindsurfCmd() *cobra.Command {
	var localOnly, globalOnly bool

	cmd := &cobra.Command{
		Use:   "windsurf",
		Short: "Install Windsurf configuration files",
		Long:  "Install local and global configuration files for Windsurf",
		Run: func(cmd *cobra.Command, args []string) {
			// Create installer instance
			installer := installer.NewWindsurfInstaller()

			// Process based on flags
			if localOnly {
				fmt.Println("Installing Windsurf local configuration file...")
				err := installer.InstallLocal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Local configuration file installation completed")
			} else if globalOnly {
				fmt.Println("Installing Windsurf global configuration file...")
				err := installer.InstallGlobal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Global configuration file installation completed")
			} else {
				fmt.Println("Installing all Windsurf configuration files...")
				err := installer.InstallAll()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("All configuration files installation completed")
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Install only local configuration file")
	cmd.Flags().BoolVarP(&globalOnly, "global", "g", false, "Install only global configuration file")

	// Make flags mutually exclusive
	cmd.MarkFlagsMutuallyExclusive("local", "global")

	return cmd
}
