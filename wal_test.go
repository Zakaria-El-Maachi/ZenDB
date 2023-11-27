package main

// import (
// 	"bytes"
// 	"io"
// 	"os"
// 	"testing"
// )

// // TestWalWrite tests the Write method of Wal.
// func TestWalWrite(t *testing.T) {
// 	wal := &Wal{
// 		file:      createTestFile(t, "test.wal"),
// 		watermark: 0,
// 	}
// 	t.Log(wal.file.Name())
// 	defer wal.file.Close()
// 	defer os.Remove(wal.file.Name())

// 	op := []byte("test operation")

// 	err := wal.Write(op)
// 	if err != nil {
// 		t.Errorf("Error writing to WAL: %v", err)
// 	}
// }

// // createTestFile creates a test file for testing purposes.
// func createTestFile(t *testing.T, fileName string) *os.File {
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		t.Fatalf("Error creating test file: %v", err)
// 	}
// 	return file
// }

// // TestWalClean tests the Clean method of Wal.
// func TestWalClean(t *testing.T) {
// 	wal := &Wal{
// 		file:      createTestFile(t, "test.wal"),
// 		watermark: 0,
// 		water:     createTestFile(t, "water_test.wal"),
// 	}

// 	defer func() {
// 		if err := wal.file.Close(); err != nil {
// 			t.Errorf("Error closing WAL file: %v", err)
// 		}
// 		if err := os.Remove(wal.file.Name()); err != nil {
// 			t.Errorf("Error removing WAL file: %v", err)
// 		}
// 	}()
// 	defer func() {
// 		if err := wal.water.Close(); err != nil {
// 			t.Errorf("Error closing water file: %v", err)
// 		}
// 		if err := os.Remove(wal.water.Name()); err != nil {
// 			t.Errorf("Error removing water file: %v", err)
// 		}
// 	}()

// 	buffer := make([]byte, 0)
// 	for j := 0; j < 100; j++ {
// 		op := []byte("operation1")
// 		buffer = append(buffer, op...)
// 		err := wal.Write(op)
// 		if err != nil {
// 			t.Errorf("Error writing to WAL - Iteration %d: %v", j, err)
// 		}
// 	}

// 	wal.watermark = 500
// 	err := wal.Clean()
// 	if err != nil {
// 		t.Errorf("Error cleaning WAL: %v", err)
// 	}

// 	buffer2 := make([]byte, 500)
// 	if _, err := io.ReadFull(wal.file, buffer2); err != nil {
// 		t.Errorf("Error reading from WAL: %v", err)
// 	}

// 	if !bytes.Equal(buffer[500:], buffer2) {
// 		t.Errorf("Error cleaning WAL : Cleaning is not properly coded")
// 	}
// }
