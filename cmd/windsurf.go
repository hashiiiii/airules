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
		Short: "Install Windsurf rules-for-ai files",
		Long:  "Install local and global rules-for-ai files for Windsurf",
		Run: func(cmd *cobra.Command, args []string) {
			// Create installer instance
			installer := installer.NewWindsurfInstaller()

			// Process based on flags
			if localOnly {
				fmt.Println("Installing Windsurf local rules-for-ai file...")
				err := installer.InstallLocal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Local rules-for-ai file installation completed")
			} else if globalOnly {
				fmt.Println("Installing Windsurf global rules-for-ai file...")
				err := installer.InstallGlobal()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("Global rules-for-ai file installation completed")
			} else {
				fmt.Println("Installing all Windsurf rules-for-ai files...")
				err := installer.InstallAll()
				if err != nil {
					fmt.Printf("Error during installation: %v\n", err)
					return
				}
				fmt.Println("All rules-for-ai files installation completed")
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Install only local rules-for-ai file")
	cmd.Flags().BoolVarP(&globalOnly, "global", "g", false, "Install only global rules-for-ai file")

	// Make flags mutually exclusive
	cmd.MarkFlagsMutuallyExclusive("local", "global")

	return cmd
}
