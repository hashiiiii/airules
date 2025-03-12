package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize airules configuration",
		Long:  "Create the configuration directory and default files",
		Run: func(cmd *cobra.Command, args []string) {
			// Get config directory
			configDir, err := config.GetConfigDir()
			if err != nil {
				fmt.Printf("Error getting config directory: %v\n", err)
				return
			}

			// Create config directory if it doesn't exist
			if err := os.MkdirAll(configDir, 0755); err != nil {
				fmt.Printf("Error creating config directory: %v\n", err)
				return
			}

			// Create default config file
			cfg := config.DefaultConfig()
			if err := config.SaveConfig(cfg); err != nil {
				fmt.Printf("Error saving default config: %v\n", err)
				return
			}

			// Create default rule file
			defaultRuleFile := filepath.Join(configDir, ".windsurfrules")
			if _, err := os.Stat(defaultRuleFile); os.IsNotExist(err) {
				// Create with default content
				defaultContent := "// Default rules-for-ai file\n// Add your rules here\n"
				if err := os.WriteFile(defaultRuleFile, []byte(defaultContent), 0644); err != nil {
					fmt.Printf("Error creating default rule file: %v\n", err)
					return
				}
			}

			fmt.Printf("Initialized airules configuration in %s\n", configDir)
			fmt.Println("Default key 'default' created with file '.windsurfrules'")
			fmt.Println("Use 'airules rules add' to add more files to keys")
		},
	}

	return cmd
}
