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

func ReadCSVFile(file *os.File, delimiter rune) ([][]string, error) {
	if filepath.Ext(file.Name()) != ".csv" {
		return nil, fmt.Errorf("File %s is not a CSV file", file.Name())
	}

	reader := csv.NewReader(file)
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

func ReadAndCleanCSVFile(file *os.File, delimiter rune) ([][]string, error) {
	if filepath.Ext(file.Name()) != ".csv" {
		return nil, fmt.Errorf("File %s is not a CSV file", file.Name())
	}

	reader := csv.NewReader(file)
	fmt.Print("Creating Reader\n")
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true
	header, err := reader.Read()
	fmt.Printf("Header: %v\n", header)

	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("File %s is empty", file.Name())
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
		fmt.Printf("All columns are empty in file %s\n", file.Name())
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

func SplitFiles(records [][]string, fileName string, fileCount int, delimiter rune) int {
	baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	length := len(records)
	lines_count := length / fileCount

	for i := range fileCount {
		index := i + 1
		fileName := fmt.Sprintf("%s-%d.csv", baseName, index)

		if index == fileCount {
			WriteSingleFile(fileName, records[i*lines_count:], delimiter)
		} else {
			WriteSingleFile(fileName, records[i*lines_count:index*lines_count], delimiter)
		}
	}
	return fileCount
}

func WriteSingleFile(fileName string, records [][]string, delimiter rune) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Comma = rune(delimiter)
	if len(headers) > 0 {
		writer.Write(headers)
	}
	err = writer.WriteAll(records)
	writer.Flush()

	return err
}
