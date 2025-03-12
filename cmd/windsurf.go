package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

// newWindsurfCmd returns the windsurf command
func newWindsurfCmd() *cobra.Command {
	var installTypeFlag string
	var keyFlag string

	cmd := &cobra.Command{
		Use:   "windsurf",
		Short: "Install Windsurf rules-for-ai files",
		Long:  "Install local and global rules-for-ai files for Windsurf",
		Run: func(cmd *cobra.Command, args []string) {
			// Create installer instance
			windsurfInstaller, err := installer.NewWindsurfInstaller()
			if err != nil {
				fmt.Printf("Error creating installer: %v\n", err)
				return
			}

			// Determine installation type based on flag
			var installType installer.InstallType
			switch installTypeFlag {
			case "local":
				installType = installer.Local
				fmt.Printf("Installing Windsurf local rules-for-ai file using key '%s'...\n", keyFlag)
			case "global":
				installType = installer.Global
				fmt.Printf("Installing Windsurf global rules-for-ai file using key '%s'...\n", keyFlag)
			default:
				installType = installer.All
				fmt.Printf("Installing all Windsurf rules-for-ai files using key '%s'...\n", keyFlag)
			}

			// Perform installation with key
			err = windsurfInstaller.InstallWithKey(installType, keyFlag)
			if err != nil {
				fmt.Printf("Error during installation: %v\n", err)
				return
			}

			fmt.Printf("%s rules-for-ai file installation completed using key '%s'\n", installTypeFlag, keyFlag)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&installTypeFlag, "type", "t", "all", "Installation type: 'local', 'global', or 'all'")
	cmd.Flags().StringVarP(&keyFlag, "key", "k", "default", "Rule file key to use")

	return cmd
}
