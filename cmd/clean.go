package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans the empty rows of the csv file",
	Long:  `Cleans the empty rows of the csv file`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string

		if len(args) == 0 {
			keymap := huh.NewDefaultKeyMap()
			keymap.FilePicker.Next = key.NewBinding(key.WithDisabled())

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewFilePicker().
						Title("Enter your csv file").
						CurrentDirectory(".").
						Value(&fileName).
						Picking(true).
						FileAllowed(true).
						DirAllowed(false).
						ShowPermissions(false).
						AllowedTypes([]string{"csv"}),

					huh.NewSelect[string]().
						Title("Delimiter").
						Value(&delimiter).
						Options(
							huh.NewOption("Comma (,)", ",").Selected(true),
							huh.NewOption("Semicolon (;)", ";"),
							huh.NewOption("Tab (  )", "\t"),
							huh.NewOption("Pipe (|)", "|"),
							huh.NewOption("Space ( )", " "),
							huh.NewOption("Colon (:)", ":"),
						),
				).
					WithHeight(20).
					Title("Split CSV"),
			).WithKeyMap(keymap)

			form.Run()
		} else {
			fileName = args[0]
		}
		selectedFile, err := os.Open(fileName)
		defer selectedFile.Close()
		if err != nil {
			log.Fatalf("error opening file %s: %v\n", fileName, err)
		}
		HandleSplit(selectedFile, 1, delimiter, true)
		fmt.Printf("Cleaned the empty rows of %s\n", fileName)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Delimiter to use")
}
