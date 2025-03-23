package tests

import (
	"csv-tools/cmd"
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func TestClean(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "testClean*.csv")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}

	csvWriter := csv.NewWriter(tmpFile)
	csvWriter.Comma = ';'
	records := [][]string{
		{"header1", "header2"},
		{"data1", "data2"},
		{"data3", "data4"},
	}
	csvWriter.WriteAll(records)

	csvWriter.Flush()
	tmpFile.Close()

	testFile, err := os.Open(tmpFile.Name())
	defer testFile.Close()
	if err != nil {
		t.Fatalf("Error opening temp file: %v", err)
	}
	// Run the clean command
	fileCount := cmd.HandleSplit(testFile, 1, ";", true)

	fmt.Printf("FileName: %s\n", testFile.Name())
	if fileCount != 1 {
		t.Errorf("Expected 1 file, got %d", fileCount)
	}
	os.Remove(testFile.Name())
	os.Remove(tmpFile.Name())
}
