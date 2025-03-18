package cmd

import (
	"fmt"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

// newInstallCmd returns the install command
func newInstallCmd() *cobra.Command {
	var editorFlag string
	var modeFlag string
	var templateFlag string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install rules-for-ai files",
		Long:  "Install rules-for-ai files for AI-powered editors like Windsurf and Cursor",
		Example: `  # Install default rules for all editors and modes
  airules install

  # Install rules only for Windsurf
  airules install --editor=windsurf

  # Install rules only for local mode
  airules install --mode=local

  # Install rules for specific editor and mode
  airules install --editor=cursor --mode=global

  # Install rules using a specific template
  airules install --template=coding_standards`,
		Run: func(cmd *cobra.Command, args []string) {
			// Parse editors
			editors := parseCommaList(editorFlag)
			if len(editors) == 0 {
				// Default to all supported editors
				editors = installer.GetSupportedEditors()
			}

			// Parse modes
			modes := parseCommaList(modeFlag)
			if len(modes) == 0 {
				// Default to all modes
				modes = []string{"local", "global"}
			}

			// Install for each editor and mode
			for _, editor := range editors {
				for _, mode := range modes {
					fmt.Printf("Installing %s rules for %s mode using template '%s'...\n", editor, mode, templateFlag)

					// Get appropriate installer
					inst, err := installer.GetInstaller(editor)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					// Determine installation type
					var installType installer.InstallType
					switch mode {
					case "local":
						installType = installer.Local
					case "global":
						installType = installer.Global
					default:
						fmt.Printf("Error: Unknown mode '%s'\n", mode)
						continue
					}

					// Install with template
					err = inst.InstallWithKey(installType, templateFlag)
					if err != nil {
						fmt.Printf("Error during installation: %v\n", err)
						continue
					}

					fmt.Printf("Successfully installed %s rules for %s mode\n", editor, mode)
				}
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor(s) to install rules for (comma-separated: windsurf,cursor or 'all')")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "", "Mode(s) to install rules for (comma-separated: local,global or 'all')")
	cmd.Flags().StringVarP(&templateFlag, "template", "t", "default", "Template key to use for installation")

	return cmd
}

// parseCommaList parses a comma-separated string into a slice of strings
func parseCommaList(s string) []string {
	if s == "" || s == "all" {
		return []string{}
	}

	var result []string
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
