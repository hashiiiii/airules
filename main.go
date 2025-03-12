package main

import (
	"fmt"
	"os"

	"github.com/hashiiiii/airules/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		// Display error message
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)

		// Display help information in the same format as --help
		rootCmd.SetArgs([]string{"--help"})
		if helpErr := rootCmd.Execute(); helpErr != nil {
			fmt.Fprintf(os.Stderr, "Error displaying help: %v\n", helpErr)
		}

		os.Exit(1)
	}
}
