package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashiiiii/airules/pkg/config"
	"github.com/hashiiiii/airules/pkg/template"
	"github.com/spf13/cobra"
)

// newTemplateCmd returns the template command
func newTemplateCmd() *cobra.Command {
	var listFlag bool
	var showFlag string
	var importFlag string
	var exportFlag string
	var outputPathFlag string
	var editorFlag string
	var modeFlag string

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage rules-for-ai templates",
		Long:  "List, show, import, and export rules-for-ai templates",
		Example: `  # List all available templates
  airules template -l

  # Show the content of a specific template
  airules template -s default -e windsurf -m global

  # Import a template from a file
  airules template -i /path/to/template.md -e windsurf -m local

  # Export a template to a file
  airules template -o default -e cursor -m global --path /path/to/output.md`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate editor flag if provided
			if editorFlag != "" {
				supportedEditors := config.GetSupportedEditors()
				valid := false
				for _, editor := range supportedEditors {
					if editorFlag == editor {
						valid = true
						break
					}
				}
				if !valid {
					fmt.Printf("Error: Invalid editor '%s'. Supported editors: %v\n",
						editorFlag, supportedEditors)
					return
				}
			} else {
				// Default to windsurf if not specified
				editorFlag = "windsurf"
			}

			// Validate mode flag if provided
			if modeFlag != "" && modeFlag != "local" && modeFlag != "global" {
				fmt.Printf("Error: Invalid mode '%s'. Must be 'local' or 'global'\n", modeFlag)
				return
			} else if modeFlag == "" {
				// Default to local if not specified
				modeFlag = "local"
			}

			// Count the number of action flags
			actionCount := 0
			if listFlag {
				actionCount++
			}
			if showFlag != "" {
				actionCount++
			}
			if importFlag != "" {
				actionCount++
			}
			if exportFlag != "" {
				actionCount++
			}

			// Ensure only one action is specified
			if actionCount == 0 {
				// Default to list if no action specified
				listFlag = true
			} else if actionCount > 1 {
				fmt.Println("Error: Only one action flag can be specified at a time")
				_ = cmd.Help()
				return
			}

			// Handle list flag
			if listFlag {
				templates, err := template.ListTemplates(editorFlag, modeFlag)
				if err != nil {
					fmt.Printf("Error listing templates: %v\n", err)
					return
				}

				fmt.Printf("Available templates for %s (%s mode):\n", editorFlag, modeFlag)
				if len(templates) == 0 {
					fmt.Println("  No templates found")
					return
				}

				for _, t := range templates {
					fmt.Printf("  - %s\n", t)
				}
				return
			}

			// Handle show flag
			if showFlag != "" {
				content, err := template.ShowTemplate(editorFlag, modeFlag, showFlag)
				if err != nil {
					fmt.Printf("Error showing template: %v\n", err)
					return
				}

				fmt.Printf("Content of template '%s' for %s (%s mode):\n\n",
					showFlag, editorFlag, modeFlag)
				fmt.Println(content)
				return
			}

			// Handle import flag
			if importFlag != "" {
				// Check if file exists
				if _, err := os.Stat(importFlag); os.IsNotExist(err) {
					fmt.Printf("Error: File '%s' does not exist\n", importFlag)
					return
				}

				// Get template name from filename if not provided
				templateName := filepath.Base(importFlag)
				// Remove extension if present
				ext := filepath.Ext(templateName)
				if ext != "" {
					templateName = templateName[:len(templateName)-len(ext)]
				}

				// Import template
				err := template.ImportTemplate(editorFlag, modeFlag, templateName, importFlag)
				if err != nil {
					fmt.Printf("Error importing template: %v\n", err)
					return
				}

				fmt.Printf("Successfully imported template '%s' for %s (%s mode)\n",
					templateName, editorFlag, modeFlag)
				return
			}

			// Handle export flag
			if exportFlag != "" {
				// Check if output path is provided
				if outputPathFlag == "" {
					fmt.Println("Error: --path is required when using -o/--export")
					return
				}

				// Export template
				err := template.ExportTemplate(editorFlag, modeFlag, exportFlag, outputPathFlag)
				if err != nil {
					fmt.Printf("Error exporting template: %v\n", err)
					return
				}

				fmt.Printf("Successfully exported template '%s' to '%s'\n",
					exportFlag, outputPathFlag)
				return
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available templates")
	cmd.Flags().StringVarP(&showFlag, "show", "s", "", "Show content of a template")
	cmd.Flags().StringVarP(&importFlag, "import", "i", "", "Import a template from a file")
	cmd.Flags().StringVarP(&exportFlag, "export", "o", "", "Export a template to a file")
	cmd.Flags().StringVar(&outputPathFlag, "path", "", "Output path for export operation")
	cmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor to operate on (default: windsurf)")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "", "Mode to operate on: 'local' or 'global' (default: local)")

	return cmd
}
