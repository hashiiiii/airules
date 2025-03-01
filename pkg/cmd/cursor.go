package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

// newCursorCmd returns the cursor command
func newCursorCmd() *cobra.Command {
	var localOnly, globalOnly bool

	cmd := &cobra.Command{
		Use:   "cursor",
		Short: "Install Cursor configuration files",
		Long:  "Install local and global configuration files for Cursor",
		Run: func(cmd *cobra.Command, args []string) {
			// Create installer instance
			installer := installer.NewCursorInstaller()

			// Process based on flags
			if localOnly {
				fmt.Println("Installing Cursor local configuration file...")
				err := installer.InstallLocal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Local configuration file installation completed")
			} else if globalOnly {
				fmt.Println("Installing Cursor global configuration file...")
				err := installer.InstallGlobal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Global configuration file installation completed")
			} else {
				fmt.Println("Installing all Cursor configuration files...")
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
