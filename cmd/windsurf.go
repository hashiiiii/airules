package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

// newWindsurfCmd returns the windsurf command
func newWindsurfCmd() *cobra.Command {
	var installTypeFlag string
	var languageFlag string

	cmd := &cobra.Command{
		Use:   "windsurf",
		Short: "Install Windsurf rules-for-ai files",
		Long:  "Install local and global rules-for-ai files for Windsurf",
		Run: func(cmd *cobra.Command, args []string) {
			var lang installer.Language
			switch languageFlag {
			case "ja", "japanese":
				lang = installer.Japanese
				fmt.Println("日本語版テンプレートを使用します...")
			default:
				lang = installer.English
				fmt.Println("Using English templates...")
			}

			// Create installer instance
			windsurfInstaller, err := installer.NewWindsurfInstaller(lang)
			if err != nil {
				fmt.Printf("Error creating installer: %v\n", err)
				return
			}

			// Determine installation type based on flag
			var installType installer.InstallType
			switch installTypeFlag {
			case "local":
				installType = installer.Local
				fmt.Println("Installing Windsurf local rules-for-ai file...")
			case "global":
				installType = installer.Global
				fmt.Println("Installing Windsurf global rules-for-ai file...")
			default:
				installType = installer.All
				fmt.Println("Installing all Windsurf rules-for-ai files...")
			}

			// Perform installation
			err = windsurfInstaller.Install(installType)
			if err != nil {
				fmt.Printf("Error during installation: %v\n", err)
				return
			}

			fmt.Printf("%s rules-for-ai file installation completed\n", installTypeFlag)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&installTypeFlag, "type", "t", "all", "Installation type: 'local', 'global', or 'all'")
	cmd.Flags().StringVarP(&languageFlag, "language", "l", "en", "Template language: 'ja' or 'en'")

	return cmd
}
