package main

import (
	"os"
	"testing"
)

// TestWalWrite tests the Write method of Wal.
func TestWalWrite(t *testing.T) {
	wal := &Wal{
		file: createTestFile(t, "test.wal"),
	}
	t.Log(wal.file.Name())
	defer wal.file.Close()

	op := []byte("test operation")

	err := wal.Write(op)
	if err != nil {
		t.Errorf("Error writing to WAL: %v", err)
	}
}

// createTestFile creates a test file for testing purposes.
func createTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	return file
}

// TestWalClean tests the Clean method of Wal.
func TestWalClean(t *testing.T) {
	wal := &Wal{
		file: createTestFile(t, "test.wal"),
	}

	defer func() {
		if err := wal.file.Close(); err != nil {
			t.Errorf("Error closing WAL file: %v", err)
		}
	}()

	buffer := make([]byte, 0)
	for j := 0; j < 100; j++ {
		op := []byte("operation1")
		buffer = append(buffer, op...)
		err := wal.Write(op)
		if err != nil {
			t.Errorf("Error writing to WAL - Iteration %d: %v", j, err)
		}
	}

	err := wal.Clean()
	if err != nil {
		t.Errorf("Error cleaning WAL: %v", err)
	}

	fileInfo, err := os.Stat(wal.file.Name())
	if err != nil {
		t.Errorf("Error getting Wal File information: %v", err)
	}

	if fileSize := fileInfo.Size(); fileSize > 0 {
		t.Errorf("Error, Wal Not truncated. New Size after Cleaning : %d", fileSize)
	}
}
