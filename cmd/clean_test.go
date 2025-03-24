package cmd_test

import (
	"csv-tools/cmd"
	"encoding/csv"
	"os"
	"testing"
)

func TestCleanPass(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "testClean*.csv")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}

	csvWriter := csv.NewWriter(tmpFile)
	records := [][]string{
		{"header1", "header2", "header3"},
		{"", "data2", "data3"},
		{"", "data4", "data5"},
	}
	csvWriter.WriteAll(records)
	tmpFile.Seek(0, 0)

	fileCount := cmd.HandleSplit(tmpFile, 1, ";", true)

	if fileCount != 1 {
		t.Errorf("Expected 1 file, got %d", fileCount)
	}

	cleanedFile, err := os.Open(tmpFile.Name()[0:len(tmpFile.Name())-4] + "-1.csv")
	if err != nil {
		t.Fatalf("Error opening cleaned file: %v", err)
	}

	tmpFile.Seek(0, 0)
	originalHeaders, _ := csv.NewReader(tmpFile).Read()
	t.Logf("Original Headers: %v\n", originalHeaders)

	cleanHeaders, _ := csv.NewReader(cleanedFile).Read()
	t.Logf("Clean Headers: %v\n", cleanHeaders)

	if originalHeaders[0] == cleanHeaders[0] {
		t.Errorf("Expected header to be cleaned, got %s", cleanHeaders[0])
	}

	t.Cleanup(func() {
		tmpFile.Close()
		cleanedFile.Close()
		os.Remove(tmpFile.Name())
		os.Remove(tmpFile.Name()[0:len(tmpFile.Name())-4] + "-1.csv")
	})
}

func TestCleanDoNotClean(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "testClean*.csv")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}

	csvWriter := csv.NewWriter(tmpFile)
	records := [][]string{
		{"header1", "header2", "header3"},
		{"", "data2", "data3"},
		{"data1", "data4", "data5"},
	}
	csvWriter.WriteAll(records)
	tmpFile.Seek(0, 0)

	fileCount := cmd.HandleSplit(tmpFile, 1, ";", true)

	if fileCount != 1 {
		t.Errorf("Expected 1 file, got %d", fileCount)
	}

	cleanedFile, err := os.Open(tmpFile.Name()[0:len(tmpFile.Name())-4] + "-1.csv")
	if err != nil {
		t.Fatalf("Error opening cleaned file: %v", err)
	}

	tmpFile.Seek(0, 0)
	originalHeaders, _ := csv.NewReader(tmpFile).Read()
	t.Logf("Original Headers: %v\n", originalHeaders)

	cleanHeaders, _ := csv.NewReader(cleanedFile).Read()
	t.Logf("Clean Headers: %v\n", cleanHeaders)

	if originalHeaders[0] != cleanHeaders[0] {
		t.Errorf("Expected headers to be the same, got %s", cleanHeaders[0])
	}

	t.Cleanup(func() {
		tmpFile.Close()
		cleanedFile.Close()
		os.Remove(tmpFile.Name())
		os.Remove(tmpFile.Name()[0:len(tmpFile.Name())-4] + "-1.csv")
	})
}
