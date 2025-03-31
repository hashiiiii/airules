package cmd

import (
	"fmt"
	"strings"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/spf13/cobra"
)

const (
	modeLocal  = "local"
	modeGlobal = "global"
)

// newInstallCmd returns the install command.
func newInstallCmd() *cobra.Command {
	var editorFlag string
	var modeFlag string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install rules-for-ai files",
		Long:  "Install rules-for-ai files for AI-powered editors like Windsurf and Cursor",
		Example: `  # Install both local and global rules for Windsurf
  airules install -e windsurf

  # Install only local rules for Cursor
  airules install -e cursor -m local

  # Install only global rules for Windsurf
  airules install -e windsurf -m global`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if editor is specified
			if editorFlag == "" {
				fmt.Println("Error: Editor must be specified using -e/--editor flag")
				fmt.Println("Supported editors:", strings.Join(installer.GetSupportedEditors(), ", "))

				return
			}

			// Check if editor is supported
			if !installer.IsEditorSupported(editorFlag) {
				fmt.Printf("Error: Unsupported editor '%s'\n", editorFlag)
				fmt.Println("Supported editors:", strings.Join(installer.GetSupportedEditors(), ", "))

				return
			}

			// Determine installation type based on mode flag
			var installType installer.InstallType
			switch modeFlag {
			case modeLocal:
				installType = installer.Local
			case modeGlobal:
				// グローバルモードが指定されたがサポートされていない場合はエラー
				if !installer.IsGlobalModeSupported(editorFlag) {
					fmt.Printf("Error: Editor '%s' does not support global mode installation through files\n", editorFlag)
					fmt.Println("Global rules for this editor must be set through the editor's settings interface")

					return
				}
				installType = installer.Global
			case "":
				// Default to both modes if not specified
				installType = installer.All
			default:
				fmt.Printf("Error: Invalid mode '%s'. Valid values are '%s' or '%s'\n", modeFlag, modeLocal, modeGlobal)

				return
			}

			// Display information about the installation
			fmt.Printf("Installing %s rules for %s editor...\n", getInstallTypeLabel(installType, editorFlag), editorFlag)

			// Install rules
			err := installer.Install(editorFlag, installType)
			if err != nil {
				fmt.Printf("Error during installation: %v\n", err)

				return
			}

			// Success message
			fmt.Printf("Successfully installed rules for %s editor\n", editorFlag)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor to install rules for (required)")
	cmd.Flags().StringVarP(
		&modeFlag,
		"mode",
		"m",
		"",
		fmt.Sprintf("Mode to install rules for: '%s', '%s', or both if not specified", modeLocal, modeGlobal),
	)
	if err := cmd.MarkFlagRequired("editor"); err != nil {
		panic(fmt.Sprintf("failed to mark 'editor' flag as required: %v", err))
	}

	return cmd
}

// getInstallTypeLabel returns a human-readable label for the install type.
func getInstallTypeLabel(installType installer.InstallType, editor string) string {
	switch installType {
	case installer.Local:
		return "local"
	case installer.Global:
		return "global"
	case installer.All:
		if installer.IsGlobalModeSupported(editor) {
			return "local and global"
		}

		return "local"
	default:
		return "unknown"
	}
}
