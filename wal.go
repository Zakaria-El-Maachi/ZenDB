package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
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

// Clean removes watermarked entries from the WAL while updating the watermark.
func (w *Wal) Clean() error {
	temp := strings.Split(w.file.Name(), ".")
	newFileName := temp[0] + "2." + temp[1]
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer func() {
		newFile.Close()
		os.Remove(newFileName)
	}()

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
		_, err = newFile.Write(buffer[:n])
		if err != nil {
			return ErrWriteFailed
		}
	}

	_, err = io.Copy(w.file, newFile)
	if err != nil {
		return ErrWriteFailed
	}

	w.watermark = 0
	buffer4 := make([]byte, 4)
	binary.LittleEndian.PutUint64(buffer4, 0)

	_, err = w.water.Seek(0, io.SeekStart)
	if err != nil {
		return ErrWriteFailed
	}

	_, err = w.water.Write(buffer4)
	if err != nil {
		return ErrWriteFailed
	}

	return nil
}
