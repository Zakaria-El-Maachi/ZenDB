package main

import (
	"io"
	"os"
	"testing"
)

var pairs = map[string]string{
	"status":    "active",
	"priority":  "high",
	"key1":      "value1",
	"injustice": "destroy",
	"soul":      "elamari",
	"zakaria":   "elmaachi",
}

// TestDecodeBytes tests the decodeBytes function.
func TestDecodeBytes(t *testing.T) {
	testFile, err := os.Open("test_file.sst")
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer testFile.Close()
	testFile.Seek(40, io.SeekStart)
	result, err := decodeBytes(testFile)
	if err != nil {
		t.Errorf("Error decoding bytes: %v", err)
	}
	if result != "injustice" {
		t.Errorf("DecodeBytes is not well implemented: %v", result)
	}

	result, err = decodeBytes(testFile)
	if err != nil {
		t.Errorf("Error decoding bytes: %v", err)
	}
	if result != "destroy" {
		t.Errorf("DecodeBytes is not well implemented: %v", result)
	}
}

// TestDecodeHeader tests the decodeHeader function.
func TestDecodeHeader(t *testing.T) {
	testFile, err := os.Open("test_file.sst")
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer testFile.Close()

	magic, entryCount, _, version, err := decodeHeader(testFile)
	if err != nil {
		t.Errorf("Error decoding header: %v", err)
	}

	if magic != MAGIC {
		t.Errorf("Error decoding header, Magic Number: %v", magic)
	}
	if entryCount != 6 {
		t.Errorf("Error decoding header,  EntryCount: %v", entryCount)
	}
	if version != 1 {
		t.Errorf("Error decoding header, Version Number: %v", version)
	}
}

// TestParse tests the parse function.
func TestParse(t *testing.T) {
	testFile, err := os.Open("test_file.sst")
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer testFile.Close()
	testFile.Seek(39, io.SeekStart)
	mem := NewMemTable()
	err = parseBody(testFile, 6, mem)
	if err != nil {
		t.Errorf("Error parsing file: %v", err)
	}

	if mem.table.size != 6 {
		t.Error("Error parsing file - Table Size")
	}
	for k, v := range pairs {
		value, err := mem.Get(k)
		if err != nil {
			t.Error("Error parsing file, probably memTable implementatio")
		}
		if value != v {
			t.Error("Error parsing file - value conflict")
		}
	}
}

// TestSearch tests the search function.
func TestSearch(t *testing.T) {
	testFile, err := os.Open("test_file.sst")
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer testFile.Close()

	for k, v := range pairs {
		value, err := Search(k, testFile)
		if err != nil {
			t.Errorf("Error searching for key: %v", k)
		}
		if value != v {
			t.Errorf("Error getting value: %v", value)
		}
	}
}
