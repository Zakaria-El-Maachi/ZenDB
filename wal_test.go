package main

import (
	"os"
	"testing"
)

// TestWalWrite tests the Write method of Wal.
func TestWalWrite(t *testing.T) {
	fileName := "test.wal"
	wal := &Wal{
		file: createTestFile(t, fileName),
	}

	defer func() {
		if err := wal.file.Close(); err != nil {
			t.Errorf("Error closing test file '%s': %v", fileName, err)
		}
		if err := os.Remove(fileName); err != nil {
			t.Errorf("Error removing test file '%s': %v", fileName, err)
		}
	}()

	op := []byte("test operation")

	err := wal.Write(op)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
}

// createTestFile creates a test file for testing purposes.
func createTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Error creating test file '%s': %v", fileName, err)
	}
	return file
}

// TestWalClean tests the Clean method of Wal.
func TestWalClean(t *testing.T) {
	fileName := "test.wal"
	wal := &Wal{
		file: createTestFile(t, fileName),
	}

	defer func() {
		if err := wal.file.Close(); err != nil {
			t.Errorf("Error closing test file '%s': %v", fileName, err)
		}
		if err := os.Remove(fileName); err != nil {
			t.Errorf("Error removing test file '%s': %v", fileName, err)
		}
	}()

	for j := 0; j < 100; j++ {
		op := []byte("operation1")
		err := wal.Write(op)
		if err != nil {
			t.Errorf("Write failed - Iteration %d: %v", j, err)
		}
	}

	err := wal.Clean()
	if err != nil {
		t.Errorf("Clean failed: %v", err)
	}

	fileInfo, err := os.Stat(wal.file.Name())
	if err != nil {
		t.Errorf("Error getting Wal File information: %v", err)
	}

	if fileSize := fileInfo.Size(); fileSize > 0 {
		t.Errorf("Error, Wal not truncated. New Size after Cleaning: %d", fileSize)
	}
}
