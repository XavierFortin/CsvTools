package csv_utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

var headers []string

func SetHeaders(h []string) {
	headers = h
}

func ReadCSVFile(fileName string, delimiter rune) ([][]string, error) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	if filepath.Ext(fileName) != ".csv" {
		return nil, fmt.Errorf("File %s is not a CSV file", fileName)
	}

	reader := csv.NewReader(csvFile)
	reader.Comma = delimiter
	reader.ReuseRecord = true
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return nil, err
	}

	return records, nil
}

func WriteCSVFile(fileName string, records [][]string, delimiter string) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Comma = rune(delimiter[0])
	if len(headers) > 0 {
		writer.Write(headers)
	}
	err = writer.WriteAll(records)
	writer.Flush()

	return err
}
