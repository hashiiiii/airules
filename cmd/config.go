package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	var listFlag bool
	var addFlag bool
	var removeFlag bool
	var removeAllFlag bool
	var modeFlag string
	var editorFlag string
	var listEditorsFlag bool

	cmd := &cobra.Command{
		Use:   "config [flags] [args...]",
		Short: "Manage AI editor rules configuration",
		Long: `Manage rules-for-ai files for AI-powered editors.

This command allows you to list, add, and remove rule files for AI-powered editors.
Each rule is identified by a key and associated with one or more files.

Rules can be configured in either 'local' (project-specific) or 'global' (all projects) mode.
Use the --editor flag to specify which editor's rules to manage.`,
		Example: `  # List all rules for the default editor in local mode
  airules config --list

  # Add a file to a rule key
  airules config --add coding_standards templates/standards.md

  # Remove a specific file from a rule key
  airules config --remove coding_standards templates/standards.md

  # Remove an entire rule key and all its files
  airules config --remove coding_standards --remove-all

  # List rules for a specific editor in global mode
  airules config --list --editor cursor --mode global

  # List supported editors
  airules config --list-editors`,
		Run: func(cmd *cobra.Command, args []string) {
			// List supported editors if requested
			if listEditorsFlag {
				supportedEditors := config.GetSupportedEditors()
				fmt.Println("Supported editors:")
				for _, editor := range supportedEditors {
					fmt.Printf("  - %s\n", editor)
				}
				return
			}

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

			// Load configuration
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			// Validate editor flag
			supportedEditors := config.GetSupportedEditors()
			editorValid := false
			for _, editor := range supportedEditors {
				if editorFlag == editor {
					editorValid = true
					break
				}
			}
			if !editorValid {
				fmt.Printf("Error: Invalid editor '%s'. Supported editors: %s\n",
					editorFlag, strings.Join(supportedEditors, ", "))
				return
			}

			// Validate mode flag
			if modeFlag != "local" && modeFlag != "global" {
				fmt.Printf("Error: Invalid mode '%s'. Must be 'local' or 'global'\n", modeFlag)
				return
			}

			// Get the appropriate map based on editor and mode
			var rulesMap map[string][]string
			editorConfig := cfg.Editors[editorFlag]
			if modeFlag == "local" {
				rulesMap = editorConfig.Local
				if rulesMap == nil {
					rulesMap = make(map[string][]string)
					editorConfig.Local = rulesMap
				}
			} else {
				rulesMap = editorConfig.Global
				if rulesMap == nil {
					rulesMap = make(map[string][]string)
					editorConfig.Global = rulesMap
				}
			}

			// List operation
			if listFlag {
				fmt.Printf("Available rule keys for %s %s mode:\n", editorFlag, modeFlag)
				if len(rulesMap) == 0 {
					fmt.Println("  No rules defined")
					return
				}
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
					fmt.Println("Error: 'config --add' requires exactly 2 arguments: <key> <file>")
					fmt.Println("Example: airules config --add coding_standards templates/standards.md")
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

				// Update the config
				if modeFlag == "local" {
					editorConfig.Local = rulesMap
				} else {
					editorConfig.Global = rulesMap
				}
				cfg.Editors[editorFlag] = editorConfig

				// Save configuration
				if err := config.SaveConfig(cfg); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					return
				}

				fmt.Printf("Added file '%s' to key '%s' in %s %s mode\n", filePath, key, editorFlag, modeFlag)
				return
			}

			// Remove operation
			if removeFlag {
				if len(args) < 1 {
					fmt.Println("Error: 'config --remove' requires at least 1 argument: <key> [file]")
					fmt.Println("Example: airules config --remove coding_standards templates/standards.md")
					fmt.Println("Example: airules config --remove coding_standards --remove-all")
					return
				}

				key := args[0]

				// Check if key exists
				if _, ok := rulesMap[key]; !ok {
					fmt.Printf("Key '%s' not found in %s %s mode\n", key, editorFlag, modeFlag)
					return
				}

				// Remove entire key if --all flag is set or no file specified
				if removeAllFlag || len(args) == 1 {
					delete(rulesMap, key)

					// Update the config
					if modeFlag == "local" {
						editorConfig.Local = rulesMap
					} else {
						editorConfig.Global = rulesMap
					}
					cfg.Editors[editorFlag] = editorConfig

					// Save configuration
					if err := config.SaveConfig(cfg); err != nil {
						fmt.Printf("Error saving config: %v\n", err)
						return
					}

					fmt.Printf("Removed key '%s' and all its files from %s %s mode\n", key, editorFlag, modeFlag)
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
					fmt.Printf("File '%s' not found in key '%s' in %s %s mode\n", filePath, key, editorFlag, modeFlag)
					return
				}

				// Update key with remaining files or delete if empty
				if len(newFiles) > 0 {
					rulesMap[key] = newFiles
				} else {
					delete(rulesMap, key)
				}

				// Update the config
				if modeFlag == "local" {
					editorConfig.Local = rulesMap
				} else {
					editorConfig.Global = rulesMap
				}
				cfg.Editors[editorFlag] = editorConfig

				// Save configuration
				if err := config.SaveConfig(cfg); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					return
				}

				fmt.Printf("Removed file '%s' from key '%s' in %s %s mode\n", filePath, key, editorFlag, modeFlag)
				return
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available rule keys and files")
	cmd.Flags().BoolVarP(&addFlag, "add", "a", false, "Add a file to a rule key (requires <key> <file> arguments)")
	cmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "Remove a file from a rule key (requires <key> [file] arguments)")
	cmd.Flags().BoolVar(&removeAllFlag, "remove-all", false, "Remove the entire key and all its files (used with --remove <key>)")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "local", "Mode to operate on: 'local' (project-specific) or 'global' (all projects)")
	cmd.Flags().StringVarP(&editorFlag, "editor", "e", "windsurf", "Editor to operate on (use --list-editors to see supported editors)")
	cmd.Flags().BoolVar(&listEditorsFlag, "list-editors", false, "List supported editors")

	return cmd
}
