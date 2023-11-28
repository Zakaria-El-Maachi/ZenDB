package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

const (
	// BufferSize represents the size of the buffer used for reading and writing.
	BufferSize = 2048
)

var (
	// ErrWriteFailed is an error when writing to the WAL fails.
	ErrWriteFailed = errors.New("write to WAL failed")
	// ErrSyncFailed is an error when syncing the WAL file fails.
	ErrSyncFailed = errors.New("syncing WAL file failed")
	// ErrReadFailed is an error when reading from the WAL file fails.
	ErrReadFailed = errors.New("read from WAL file failed")
)

// Wal represents the Write-Ahead Log.
type Wal struct {
	watermark int64
	file      *os.File
	water     *os.File
	meta      *os.File
	notifier  chan bool
}

// Write appends the given operation to the WAL and flushes immediately.
func (w *Wal) Write(op []byte) error {
	if _, err := w.file.Write(op); err != nil {
		return ErrWriteFailed
	}
	if err := w.file.Sync(); err != nil {
		return ErrSyncFailed
	}
	return nil
}

// Clean removes watermarked entries from the WAL while updating the watermark, in an atomic way
func (w *Wal) Clean() error {
	log.Println("Cleaning Wal Started")

	// Get the current file name
	fileName := w.file.Name()

	// Close the file to release resources (but don't close the reference)
	w.file.Close()
	// Open the file with the same name
	var err error
	w.file, err = os.OpenFile(fileName, os.O_RDWR, FilePermission)
	if err != nil {
		return ErrReadFailed
	}

	// Create a temporary file
	tempFile, err := os.Create("tempfile.wal")
	if err != nil {
		return ErrWriteFailed
	}

	buffer := make([]byte, BufferSize)
	_, err = w.file.Seek(w.watermark, io.SeekStart)
	if err != nil {
		return ErrReadFailed
	}

	for {
		n, err := w.file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return ErrReadFailed
		}
		_, err = tempFile.Write(buffer[:n])
		if err != nil {
			return ErrWriteFailed
		}
	}

	tempFile.Close()
	w.file.Close()

	// Rename the temporary file to the original file
	err = os.Rename("tempfile.wal", fileName)
	if err != nil {
		return err
	}
	w.watermark = 0

	// Update watermark in water file
	buffer8 := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer8, 0)

	_, err = w.water.Seek(0, io.SeekStart)
	if err != nil {
		return ErrWriteFailed
	}

	_, err = w.water.Write(buffer8)
	if err != nil {
		return ErrWriteFailed
	}

	return nil
}
