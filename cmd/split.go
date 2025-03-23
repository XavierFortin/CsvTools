package cmd

import (
	"csv-tools/csv_utils"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	num_files int
	delimiter string
	clean     bool
)

func HandleSplit(file *os.File, fileCount int, delim string, clean bool) int {
	var records [][]string

	var err error
	if clean {
		records, err = csv_utils.ReadAndCleanCSVFile(file, rune(delim[0]))
	} else {
		records, err = csv_utils.ReadCSVFile(file, rune(delim[0]))
	}

	if err != nil {
		log.Fatalf("error opening file %s: %v\n", file.Name(), err)
	}

	if fileCount >= 25 {
		confirmValue := false
		huh.NewConfirm().
			Title(fmt.Sprintf("This will create %d files. Are you sure you want to continue?", fileCount)).
			Value(&confirmValue).
			Run()
		if !confirmValue {
			return 0
		}
	}

	csv_utils.SetHeaders(records[0])
	records = records[1:]
	csv_utils.SplitFiles(records, file.Name(), fileCount, rune(delim[0]))

	return fileCount
}

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split input_file",
	Short: "Splits a CSV file into multiple files",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		var filesCount string
		clean = true

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

					huh.NewInput().
						Title(fmt.Sprintf("Enter number of files to split into")).
						Validate(func(s string) error {
							if s != "" {
								_, err := strconv.Atoi(s)
								if err != nil {
									return fmt.Errorf("Invalid number of files")
								}
							}
							return nil
						}).
						Value(&filesCount),

					huh.NewConfirm().
						Title("Remove empty columns").Accessor(huh.NewPointerAccessor(&clean)).
						Description("Remove columns that are empty (contains only a header)").
						Value(&clean),
				).
					WithHeight(20).
					Title("Split CSV"),
			).WithKeyMap(keymap)

			form.Run()
		} else {
			fileName = args[0]
		}
		num_files, _ := strconv.Atoi(filesCount)
		if num_files == 0 {
			fmt.Println("Number of files to split into is required")
			return
		}
		selectedFile, err := os.Open(fileName)
		defer selectedFile.Close()
		if err != nil {
			log.Fatalf("error opening file %s: %v\n", fileName, err)
		}
		generated_files := HandleSplit(selectedFile, num_files, delimiter, clean)
		fmt.Printf("Split %s into %d files\n", fileName, generated_files)
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().SortFlags = false
	splitCmd.Flags().IntVarP(&num_files, "files", "f", 0, "Number of files to split into")
	splitCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Delimiter to use")
	splitCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Remove empty columns")
}
