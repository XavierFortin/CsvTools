package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI to the latest version",
	Long: `Check GitHub for a newer version of the CLI and update if available.
This will replace the current binary with the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		doSelfUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func doSelfUpdate() {
	// Your GitHub repository in the format "owner/repo"
	repo := "yourUsername/yourRepoName"

	// Check for the latest version
	fmt.Println("Checking for updates...")
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error finding latest version:", err)
		os.Exit(1)
	}

	if !found {
		fmt.Println("No release found on GitHub")
		return
	}

	// Remove the 'v' prefix if it exists for version comparison
	latestVersion := latest.Version.String()

	// Compare versions
	if Version == latestVersion {
		fmt.Printf("Current version (%s) is the latest\n", Version)
		return
	}

	// We have a newer version, ask the user if they want to update
	fmt.Printf("New version %s found\n", latest.Version)
	fmt.Printf("Release notes:\n%s\n", latest.ReleaseNotes)
	fmt.Print("Do you want to update? (y/N): ")

	input := ""
	fmt.Scanln(&input)
	if strings.ToLower(input) != "y" {
		fmt.Println("Update canceled")
		return
	}

	// Perform the update
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not locate executable path:", err)
		os.Exit(1)
	}

	fmt.Println("Downloading update...")
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		fmt.Fprintln(os.Stderr, "Error updating binary:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated to version %s\n", latest.Version)
}
