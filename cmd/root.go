package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "DEV"

var showVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csv_tools",
	Short: "A suite of tools for working with CSV files",

	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("CsvTools version: %s\n", Version)
			return
		}

	},
}

func Execute() {

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
