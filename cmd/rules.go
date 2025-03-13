package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/spf13/cobra"
)

func newRulesCmd() *cobra.Command {
	var listFlag bool
	var addFlag bool
	var removeFlag bool
	var removeAllFlag bool
	var modeFlag string

	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage rules-for-ai files",
		Long:  "List, add, and remove rules-for-ai files",
		Run: func(cmd *cobra.Command, args []string) {
			// Error if multiple operation flags are specified
			flagCount := 0
			if listFlag {
				flagCount++
			}
			if addFlag {
				flagCount++
			}
			if removeFlag {
				flagCount++
			}

			if flagCount > 1 {
				fmt.Println("Error: Only one operation flag can be specified at a time")
				_ = cmd.Help()
				return
			}

			// Default to list if no flags specified
			if flagCount == 0 {
				listFlag = true
			}

			// Validate mode flag
			if modeFlag != "local" && modeFlag != "global" {
				fmt.Printf("Error: Invalid mode '%s'. Must be 'local' or 'global'\n", modeFlag)
				return
			}

			// Load configuration
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			// Get the appropriate map based on mode
			var rulesMap map[string][]string
			if modeFlag == "local" {
				rulesMap = cfg.Windsurf.Local
			} else {
				rulesMap = cfg.Windsurf.Global
			}

			// List operation
			if listFlag {
				fmt.Printf("Available rule keys for windsurf %s mode:\n", modeFlag)
				for key, files := range rulesMap {
					fmt.Printf("  %s:\n", key)
					for _, file := range files {
						fmt.Printf("    - %s\n", file)
					}
				}
				return
			}

			// Add operation
			if addFlag {
				if len(args) != 2 {
					fmt.Println("Error: 'rules --add' requires exactly 2 arguments: <key> <file>")
					return
				}

				key := args[0]
				filePath := args[1]

				// Get config directory
				configDir, err := config.GetConfigDir()
				if err != nil {
					fmt.Printf("Error getting config directory: %v\n", err)
					return
				}

				// Check if file exists
				absFilePath := filepath.Join(configDir, filePath)
				if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
					fmt.Printf("Warning: File '%s' does not exist\n", absFilePath)
				}

				// Add file to key
				if _, ok := rulesMap[key]; !ok {
					rulesMap[key] = []string{}
				}

				// Check if file already exists in the key
				for _, f := range rulesMap[key] {
					if f == filePath {
						fmt.Printf("File '%s' is already in key '%s'\n", filePath, key)
						return
					}
				}

				// Add file to key
				rulesMap[key] = append(rulesMap[key], filePath)

				// Update the config based on mode
				if modeFlag == "local" {
					cfg.Windsurf.Local = rulesMap
				} else {
					cfg.Windsurf.Global = rulesMap
				}

				// Save configuration
				if err := config.SaveConfig(cfg); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					return
				}

				fmt.Printf("Added file '%s' to key '%s' in %s mode\n", filePath, key, modeFlag)
				return
			}

			// Remove operation
			if removeFlag {
				if len(args) < 1 {
					fmt.Println("Error: 'rules --remove' requires at least 1 argument: <key> [file]")
					return
				}

				key := args[0]

				// Check if key exists
				if _, ok := rulesMap[key]; !ok {
					fmt.Printf("Key '%s' not found in %s mode\n", key, modeFlag)
					return
				}

				// Remove entire key if --all flag is set or no file specified
				if removeAllFlag || len(args) == 1 {
					delete(rulesMap, key)

					// Update the config based on mode
					if modeFlag == "local" {
						cfg.Windsurf.Local = rulesMap
					} else {
						cfg.Windsurf.Global = rulesMap
					}

					// Save configuration
					if err := config.SaveConfig(cfg); err != nil {
						fmt.Printf("Error saving config: %v\n", err)
						return
					}

					fmt.Printf("Removed key '%s' and all its files from %s mode\n", key, modeFlag)
					return
				}

				// Remove specific file from key
				filePath := args[1]
				var newFiles []string
				fileFound := false

				for _, f := range rulesMap[key] {
					if f != filePath {
						newFiles = append(newFiles, f)
					} else {
						fileFound = true
					}
				}

				if !fileFound {
					fmt.Printf("File '%s' not found in key '%s' in %s mode\n", filePath, key, modeFlag)
					return
				}

				// Update key with remaining files or delete if empty
				if len(newFiles) > 0 {
					rulesMap[key] = newFiles
				} else {
					delete(rulesMap, key)
				}

				// Update the config based on mode
				if modeFlag == "local" {
					cfg.Windsurf.Local = rulesMap
				} else {
					cfg.Windsurf.Global = rulesMap
				}

				// Save configuration
				if err := config.SaveConfig(cfg); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					return
				}

				fmt.Printf("Removed file '%s' from key '%s' in %s mode\n", filePath, key, modeFlag)
				return
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available rule keys and files")
	cmd.Flags().BoolVarP(&addFlag, "add", "a", false, "Add a file to a rule key")
	cmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "Remove a file from a rule key")
	cmd.Flags().BoolVar(&removeAllFlag, "all", false, "Remove the entire key and all its files (used with --remove)")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "local", "Mode to operate on: 'local' or 'global'")

	return cmd
}
