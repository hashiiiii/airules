package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashiiiii/airules/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		// Display error message
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)

		// For unknown commands, show a special message with available commands
		if strings.Contains(err.Error(), "unknown command") {
			fmt.Fprintf(os.Stderr, "The specified command does not exist. Please check the available commands below:\n\n")
			// Display help information
			rootCmd.Help()
		} else {
			// For other errors, suggest using help
			fmt.Fprintf(os.Stderr, "For usage information, run 'airules --help'\n")
		}

		os.Exit(1)
	}
}
