/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"csv-tools/csv_utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	num_files, num_lines int
	delimiter            string
	clean                bool
)

func splitByFiles(records [][]string, fileName string, files_count int) {
	length := len(records)
	records_per_file := length / files_count

	splitByLines(records, fileName, records_per_file)
}

func splitByLines(records [][]string, fileName string, lines_count int) {
	length := len(records)

	if length < lines_count {
		fmt.Printf("File %s has less than %d lines\n", fileName, lines_count)
	}
	files_count := length / lines_count

	if files_count >= 25 {
		confirmValue := false
		huh.NewConfirm().
			Title(fmt.Sprintf("This will create %d files. Are you sure you want to continue?", files_count)).
			Value(&confirmValue).
			Run()

		if !confirmValue {
			return
		}
	}

	csv_utils.SetHeaders(records[0])
	records = records[1:]

	baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	for i := 0; i < files_count; i++ {
		index := i + 1
		fileName := fmt.Sprintf("%s-%d.csv", baseName, index)

		if index == files_count {
			csv_utils.WriteCSVFile(fileName, records[i*lines_count:], delimiter)
		} else {
			csv_utils.WriteCSVFile(fileName, records[i*lines_count:index*lines_count], delimiter)
		}
	}

	fmt.Printf("Split %s into %d files\n", fileName, files_count)

}

func HandleSplit(fileName string, lines_count int, file_count int, delimiter string, clean bool) {
	var records [][]string
	var err error
	if clean {
		records, err = csv_utils.ReadAndCleanCSVFile(fileName, rune(delimiter[0]))
	} else {
		records, err = csv_utils.ReadCSVFile(fileName, rune(delimiter[0]))
	}

	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", fileName, err)
		return
	}

	if file_count != 0 {
		splitByFiles(records, fileName, file_count)
	} else if lines_count != 0 {
		splitByLines(records, fileName, lines_count)
	}
}

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split input_file",
	Short: "Splits a CSV file into multiple files",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		clean = true

		if len(args) == 0 {
			var split_type string
			var number_of_type string
			var delimiter string

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

					huh.NewSelect[string]().
						Title("Split by number of lines or files").
						Value(&split_type).
						Options(
							huh.NewOption("Lines", "lines"),
							huh.NewOption("Files", "files"),
						),

					huh.NewInput().
						TitleFunc(func() string {
							return fmt.Sprintf("Enter number of %s to split into", split_type)
						}, &split_type).
						Validate(func(s string) error {
							if s != "" {
								_, err := strconv.Atoi(s)
								if err != nil {
									return fmt.Errorf("Invalid number of lines")
								}
							}
							return nil
						}).
						Value(&number_of_type),

					huh.NewConfirm().
						Title("Remove empty columns").Accessor(huh.NewPointerAccessor(&clean)).
						Description("Remove columns that are empty (contains only a header)").
						Value(&clean),
				).
					WithHeight(20).
					Title("Split CSV"),
			).WithKeyMap(keymap)

			form.Run()

			if split_type == "lines" {
				num_lines, _ = strconv.Atoi(number_of_type)
			} else {
				num_files, _ = strconv.Atoi(number_of_type)
			}

		} else {
			fileName = args[0]
		}

		if num_files != 0 && num_lines != 0 {
			fmt.Println("Cannot specify both files and lines")
			return
		}
		HandleSplit(fileName, num_lines, num_files, delimiter, clean)
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().SortFlags = false
	splitCmd.Flags().IntVarP(&num_files, "files", "f", 0, "Number of files to split into")
	splitCmd.Flags().IntVarP(&num_lines, "lines", "l", 0, "Number of lines per file")
	splitCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Delimiter to use")
	splitCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Remove empty columns")

	if len(os.Args) == 3 {
		splitCmd.MarkFlagsOneRequired("files", "lines")
	}
}
