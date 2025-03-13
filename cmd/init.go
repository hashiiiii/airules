package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/otiai10/copy"
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
				fmt.Printf("Failed to get config directory: %v\n", err)
				return
			}

			// Create config directory if it doesn't exist
			if err := os.MkdirAll(configDir, 0755); err != nil {
				fmt.Printf("Failed to create config directory: %v\n", err)
				return
			}

			// Get repository root directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				fmt.Println("Failed to get current file path")
				return
			}

			// Go up two directories from cmd/init.go to reach repository root
			repoRoot := filepath.Dir(filepath.Dir(filename))

			// Source templates directory (in repository root)
			srcTemplatesDir := filepath.Join(repoRoot, "templates")

			// Check if templates directory exists
			if _, err := os.Stat(srcTemplatesDir); os.IsNotExist(err) {
				fmt.Printf("'templates' directory not found at %s\n", srcTemplatesDir)
				return
			}

			// Destination templates directory
			destTemplatesDir := filepath.Join(configDir, "templates")

			// Copy templates
			fmt.Printf("Copying templates from %s to %s\n", srcTemplatesDir, destTemplatesDir)
			if err := copy.Copy(srcTemplatesDir, destTemplatesDir); err != nil {
				fmt.Printf("Failed to copy templates: %v\n", err)
				return
			}
			fmt.Println("Templates copied successfully.")

			// Create or update config file
			configFile := filepath.Join(configDir, "config.toml")
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				// Create default config with settings for both editors and modes
				cfg := config.GetDefaultConfig()
				if err := config.SaveConfig(cfg); err != nil {
					fmt.Printf("Failed to save default config: %v\n", err)
					return
				}
				fmt.Println("Created configuration file with default settings.")
			} else {
				fmt.Println("Configuration file already exists, not overwriting.")
				fmt.Println("If you want to see the current configuration, check the file at:")
				fmt.Println(configFile)
			}
		},
	}

	return cmd
}
