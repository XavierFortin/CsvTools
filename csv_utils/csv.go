package csv_utils

import (
	"encoding/csv"
	"fmt"
	"io"
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

func ReadAndCleanCSVFile(fileName string, delimiter rune) ([][]string, error) {
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
	reader.TrimLeadingSpace = true
	header, err := reader.Read()

	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("File %s is empty", fileName)
		} else {
			return nil, err
		}
	}

	numColumns := len(header)
	isEmptyColumns := make([]bool, numColumns)

	for i := range isEmptyColumns {
		isEmptyColumns[i] = true
	}

	oldRecords := make([][]string, 0)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error reading CSV file: %v", err)
		}

		for i, cell := range row {
			if i < numColumns && cell != "" {
				isEmptyColumns[i] = false
			}
		}
		oldRecords = append(oldRecords, row)
	}

	outputHeader := make([]string, 0)
	columnIndexMap := make(map[int]int)
	newColumnIndex := 0

	for i, empty := range isEmptyColumns {
		if !empty {
			outputHeader = append(outputHeader, header[i])
			columnIndexMap[i] = newColumnIndex
			newColumnIndex++
		}
	}

	newRecords := make([][]string, 0)

	if len(outputHeader) == 0 {
		fmt.Printf("All columns are empty in file %s\n", fileName)
	}

	newRecords = append(newRecords, outputHeader)

	for _, oldRecord := range oldRecords {
		newRecord := make([]string, len(outputHeader))
		for i, cell := range oldRecord {
			if newIndex, ok := columnIndexMap[i]; ok {
				newRecord[newIndex] = cell
			}
		}
		newRecords = append(newRecords, newRecord)
	}

	return newRecords, nil
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
