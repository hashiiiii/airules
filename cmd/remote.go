package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/hashiiiii/airules/pkg/remote"
	"github.com/spf13/cobra"
)

// newRemoteCmd returns the remote command
func newRemoteCmd() *cobra.Command {
	var installTypeFlag string
	var listFlag bool
	var installFlag string

	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Install rule sets from remote repositories",
		Long:  "Install rule sets from remote repositories like awesome-cursorrules",
		Run: func(cmd *cobra.Command, args []string) {
			// Create a context
			ctx := context.Background()

			// Create a fetcher
			fetcher := remote.NewGitHubFetcher(nil)

			// Create an installer
			remoteInstaller := installer.NewRemoteInstaller(fetcher)

			// Determine the installation type
			var installType installer.InstallType
			switch installTypeFlag {
			case "local":
				installType = installer.Local
			case "global":
				installType = installer.Global
			default:
				installType = installer.Local
			}

			// If list flag is set, list available rule sets
			if listFlag {
				fmt.Println("Listing available rule sets...")
				ruleSets, err := remoteInstaller.ListRuleSets(ctx)
				if err != nil {
					fmt.Printf("Error listing rule sets: %v\n", err)
					return
				}

				// Print the rule sets in a table
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "NAME\tTYPE\tDESCRIPTION\tURL")
				for _, ruleSet := range ruleSets {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ruleSet.Name, ruleSet.Type, ruleSet.Description, ruleSet.URL)
				}
				w.Flush()
				return
			}

			// If install flag is set, install the specified rule set
			if installFlag != "" {
				fmt.Printf("Installing rule set %s...\n", installFlag)
				err := remoteInstaller.InstallRuleSetByName(ctx, installFlag, installType)
				if err != nil {
					fmt.Printf("Error installing rule set: %v\n", err)
					return
				}
				fmt.Printf("Rule set %s installed successfully\n", installFlag)
				return
			}

			// If no flags are set, show an interactive menu
			fmt.Println("Fetching available rule sets...")
			ruleSets, err := remoteInstaller.ListRuleSets(ctx)
			if err != nil {
				fmt.Printf("Error fetching rule sets: %v\n", err)
				return
			}

			// Show an interactive menu
			fmt.Println("Available rule sets:")
			for i, ruleSet := range ruleSets {
				fmt.Printf("%d. %s (%s)\n", i+1, ruleSet.Name, ruleSet.Type)
			}

			// Prompt for selection
			var selection int
			fmt.Print("Enter the number of the rule set to install (or 0 to cancel): ")
			_, err = fmt.Scanln(&selection)
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return
			}

			// Check if the selection is valid
			if selection <= 0 || selection > len(ruleSets) {
				fmt.Println("Installation cancelled")
				return
			}

			// Install the selected rule set
			selectedRuleSet := ruleSets[selection-1]
			fmt.Printf("Installing rule set %s...\n", selectedRuleSet.Name)
			err = remoteInstaller.InstallRuleSet(ctx, selectedRuleSet, installType)
			if err != nil {
				fmt.Printf("Error installing rule set: %v\n", err)
				return
			}
			fmt.Printf("Rule set %s installed successfully\n", selectedRuleSet.Name)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&installTypeFlag, "type", "t", "local", "Installation type: 'local' or 'global'")
	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available rule sets")
	cmd.Flags().StringVarP(&installFlag, "install", "i", "", "Install the specified rule set")

	return cmd
}
