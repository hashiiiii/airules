package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/hashiiiii/airules/pkg/installer"
	"github.com/hashiiiii/airules/pkg/remote"
	"github.com/spf13/cobra"
)

// newWindsurfCmd returns the windsurf command
func newWindsurfCmd() *cobra.Command {
	var installTypeFlag string
	var languageFlag string
	var installFlag string
	var listFlag bool

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

			// Determine installation type based on flag
			var installType installer.InstallType
			switch installTypeFlag {
			case "local", "l":
				installType = installer.Local
				fmt.Println("Installing Windsurf local rules-for-ai file...")
			case "global", "g":
				installType = installer.Global
				fmt.Println("Installing Windsurf global rules-for-ai file...")
			default:
				installType = installer.All
				fmt.Println("Installing all Windsurf rules-for-ai files...")
			}

			// If list flag is set, list available rule sets
			if listFlag {
				listWindsurfRuleSets()
				return
			}

			// If install flag is set, install the specified rule set
			if installFlag != "" {
				installWindsurfRuleSet(installFlag, installType)
				return
			}

			// If no flags are set, show an interactive menu
			if !listFlag && installFlag == "" {
				// First try to install from remote repository
				if installWindsurfFromRemote(installType) {
					return
				}

				// If remote installation fails or is cancelled, fall back to local templates
				installWindsurfFromLocalTemplates(lang, installType)
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&installTypeFlag, "type", "t", "all", "Installation type: 'local', 'global', or 'all' (default)")
	cmd.Flags().StringVarP(&languageFlag, "language", "l", "en", "Template language: 'ja' or 'en' (default)")
	cmd.Flags().StringVarP(&installFlag, "install", "i", "", "Install the specified rule set by name or ID")
	cmd.Flags().BoolVar(&listFlag, "list", false, "List available rule sets")

	return cmd
}

// installWindsurfFromLocalTemplates installs Windsurf rules from local templates
func installWindsurfFromLocalTemplates(lang installer.Language, installType installer.InstallType) {
	// Create installer instance
	windsurfInstaller, err := installer.NewWindsurfInstaller(lang)
	if err != nil {
		fmt.Printf("Error creating installer: %v\n", err)
		return
	}

	// Perform installation
	err = windsurfInstaller.Install(installType)
	if err != nil {
		fmt.Printf("Error during installation: %v\n", err)
		return
	}

	fmt.Printf("%s rules-for-ai file installation completed\n", installType.String())
}

// installWindsurfFromRemote attempts to install Windsurf rules from remote repository
// Returns true if installation was successful or explicitly cancelled
func installWindsurfFromRemote(installType installer.InstallType) bool {
	// Create a context
	ctx := context.Background()

	// Create a fetcher
	fetcher := remote.NewGitHubFetcher(nil)

	// Create an installer
	remoteInstaller := installer.NewRemoteInstaller(fetcher)

	// Fetch available rule sets
	fmt.Println("Fetching available rule sets...")
	ruleSets, err := remoteInstaller.ListRuleSets(ctx)
	if err != nil {
		fmt.Printf("Error fetching rule sets: %v\n", err)
		fmt.Println("Falling back to local templates...")
		return false
	}

	// Filter for Windsurf rule sets
	var windsurfRuleSets []remote.RuleSet
	for _, ruleSet := range ruleSets {
		if strings.ToLower(ruleSet.Type) == "windsurf" {
			windsurfRuleSets = append(windsurfRuleSets, ruleSet)
		}
	}

	if len(windsurfRuleSets) == 0 {
		fmt.Println("No Windsurf rule sets found. Falling back to local templates...")
		return false
	}

	// Show an interactive menu
	fmt.Println("Available Windsurf rule sets:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
	for i, ruleSet := range windsurfRuleSets {
		// Get a shorter name by removing common suffixes
		shortName := ruleSet.Name
		shortName = strings.TrimSuffix(shortName, "-windsurfrules-prompt-file")
		shortName = strings.TrimSuffix(shortName, "-windsurfrules-prompt")
		shortName = strings.TrimSuffix(shortName, "-windsurfrules")

		fmt.Fprintf(w, "%d\t%s\t%s\n", i+1, shortName, ruleSet.Description)
	}
	w.Flush()

	// Prompt for selection
	var selection int
	fmt.Print("\nEnter the number of the rule set to install (or 0 to cancel): ")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return false
	}

	selection, err = strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return false
	}

	// Check if the selection is valid
	if selection <= 0 || selection > len(windsurfRuleSets) {
		fmt.Println("Installation cancelled")
		return true
	}

	// Install the selected rule set
	selectedRuleSet := windsurfRuleSets[selection-1]
	fmt.Printf("Installing rule set %s...\n", selectedRuleSet.Name)
	err = remoteInstaller.InstallRuleSet(ctx, selectedRuleSet, installType)
	if err != nil {
		fmt.Printf("Error installing rule set: %v\n", err)
		fmt.Println("Falling back to local templates...")
		return false
	}

	fmt.Printf("Rule set %s installed successfully\n", selectedRuleSet.Name)
	return true
}

// listWindsurfRuleSets lists available Windsurf rule sets
func listWindsurfRuleSets() {
	// Create a context
	ctx := context.Background()

	// Create a fetcher
	fetcher := remote.NewGitHubFetcher(nil)

	// Create an installer
	remoteInstaller := installer.NewRemoteInstaller(fetcher)

	// Fetch available rule sets
	fmt.Println("Listing available Windsurf rule sets...")
	ruleSets, err := remoteInstaller.ListRuleSets(ctx)
	if err != nil {
		fmt.Printf("Error listing rule sets: %v\n", err)
		return
	}

	// Filter for Windsurf rule sets
	var windsurfRuleSets []remote.RuleSet
	for _, ruleSet := range ruleSets {
		if strings.ToLower(ruleSet.Type) == "windsurf" {
			windsurfRuleSets = append(windsurfRuleSets, ruleSet)
		}
	}

	// Print the rule sets in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tURL")
	for i, ruleSet := range windsurfRuleSets {
		// Get a shorter name by removing common suffixes
		shortName := ruleSet.Name
		shortName = strings.TrimSuffix(shortName, "-windsurfrules-prompt-file")
		shortName = strings.TrimSuffix(shortName, "-windsurfrules-prompt")
		shortName = strings.TrimSuffix(shortName, "-windsurfrules")

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, shortName, ruleSet.Description, ruleSet.URL)
	}
	w.Flush()
}

// installWindsurfRuleSet installs a specific Windsurf rule set
func installWindsurfRuleSet(nameOrID string, installType installer.InstallType) {
	// Create a context
	ctx := context.Background()

	// Create a fetcher
	fetcher := remote.NewGitHubFetcher(nil)

	// Create an installer
	remoteInstaller := installer.NewRemoteInstaller(fetcher)

	// Check if the input is an ID
	id, err := strconv.Atoi(nameOrID)
	if err == nil {
		// Fetch available rule sets
		ruleSets, err := remoteInstaller.ListRuleSets(ctx)
		if err != nil {
			fmt.Printf("Error fetching rule sets: %v\n", err)
			return
		}

		// Filter for Windsurf rule sets
		var windsurfRuleSets []remote.RuleSet
		for _, ruleSet := range ruleSets {
			if strings.ToLower(ruleSet.Type) == "windsurf" {
				windsurfRuleSets = append(windsurfRuleSets, ruleSet)
			}
		}

		// Check if the ID is valid
		if id <= 0 || id > len(windsurfRuleSets) {
			fmt.Printf("Invalid rule set ID: %d\n", id)
			return
		}

		// Install the rule set
		selectedRuleSet := windsurfRuleSets[id-1]
		fmt.Printf("Installing rule set %s...\n", selectedRuleSet.Name)
		err = remoteInstaller.InstallRuleSet(ctx, selectedRuleSet, installType)
		if err != nil {
			fmt.Printf("Error installing rule set: %v\n", err)
			return
		}

		fmt.Printf("Rule set %s installed successfully\n", selectedRuleSet.Name)
		return
	}

	// If not an ID, try to find by name
	// Fetch available rule sets
	ruleSets, err := remoteInstaller.ListRuleSets(ctx)
	if err != nil {
		fmt.Printf("Error fetching rule sets: %v\n", err)
		return
	}

	// Filter for Windsurf rule sets and find matches
	var matchingRuleSets []remote.RuleSet
	for _, ruleSet := range ruleSets {
		if strings.ToLower(ruleSet.Type) != "windsurf" {
			continue
		}

		// Check for exact match
		if ruleSet.Name == nameOrID {
			matchingRuleSets = []remote.RuleSet{ruleSet}
			break
		}

		// Check for prefix match
		if strings.HasPrefix(ruleSet.Name, nameOrID) {
			matchingRuleSets = append(matchingRuleSets, ruleSet)
		}
	}

	if len(matchingRuleSets) == 0 {
		fmt.Printf("No rule sets found matching '%s'\n", nameOrID)
		return
	}

	if len(matchingRuleSets) == 1 {
		// Install the rule set
		selectedRuleSet := matchingRuleSets[0]
		fmt.Printf("Installing rule set %s...\n", selectedRuleSet.Name)
		err = remoteInstaller.InstallRuleSet(ctx, selectedRuleSet, installType)
		if err != nil {
			fmt.Printf("Error installing rule set: %v\n", err)
			return
		}

		fmt.Printf("Rule set %s installed successfully\n", selectedRuleSet.Name)
		return
	}

	// Multiple matches found, prompt for selection
	fmt.Printf("Multiple rule sets match '%s':\n", nameOrID)
	for i, ruleSet := range matchingRuleSets {
		fmt.Printf("%d. %s\n", i+1, ruleSet.Name)
	}

	// Prompt for selection
	var selection int
	fmt.Print("\nEnter the number of the rule set to install (or 0 to cancel): ")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	selection, err = strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	// Check if the selection is valid
	if selection <= 0 || selection > len(matchingRuleSets) {
		fmt.Println("Installation cancelled")
		return
	}

	// Install the selected rule set
	selectedRuleSet := matchingRuleSets[selection-1]
	fmt.Printf("Installing rule set %s...\n", selectedRuleSet.Name)
	err = remoteInstaller.InstallRuleSet(ctx, selectedRuleSet, installType)
	if err != nil {
		fmt.Printf("Error installing rule set: %v\n", err)
		return
	}

	fmt.Printf("Rule set %s installed successfully\n", selectedRuleSet.Name)
}
