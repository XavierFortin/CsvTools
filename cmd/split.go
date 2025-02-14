/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"csv-tools/csv_utils"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
)

var (
	num_files, num_lines int
	delimiter            string
	hasHeader            bool
	clean                bool
)

// Split the CSV file into multiple files specified by num_files
func splitByFiles(records [][]string, fileName string, files_count int) {
	length := len(records)
	records_per_file := length / files_count

	splitByLines(records, fileName, records_per_file)
}

func splitByLines(records [][]string, fileName string, lines_count int) {
	length := len(records)

	if length < num_lines {
		fmt.Printf("File %s has less than %d lines\n", fileName, num_lines)
	}
	files_count := length / lines_count

	if hasHeader {
		csv_utils.SetHeaders(records[0])
		records = records[1:]
	}

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

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split input_file",
	Short: "Splits a CSV file into multiple files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		// Split the CSV file
		if num_files != 0 && num_lines != 0 {
			fmt.Println("Cannot specify both files and lines")
			return
		}

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

		if num_files != 0 {
			splitByFiles(records, fileName, num_files)
		} else if num_lines != 0 {
			splitByLines(records, fileName, num_lines)
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().SortFlags = false
	splitCmd.Flags().IntVarP(&num_files, "files", "f", 0, "Number of files to split into")
	splitCmd.Flags().IntVarP(&num_lines, "lines", "l", 0, "Number of lines per file")
	splitCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Delimiter to use")
	splitCmd.Flags().BoolVarP(&hasHeader, "header", "H", true, "Recopy the headers for each file")
	splitCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Remove empty columns")

	splitCmd.MarkFlagsOneRequired("files", "lines")
}
