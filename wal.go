package main

import (
	"errors"
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
	file *os.File
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

// RecordSet records a 'set' operation in the WAL.
func (w *Wal) RecordSet(key, value string) error {
	op := append([]byte("s"), encodeString(key)...)
	op = append(op, encodeString(value)...)
	return w.Write(op)
}

// RecordDel records a 'delete' operation in the WAL.
func (w *Wal) RecordDel(key string) error {
	op := append([]byte("d"), encodeString(key)...)
	return w.Write(op)
}

// Clean removes watermarked entries from the WAL while updating the watermark, in an atomic way
func (w *Wal) Clean() error {
	// Close the current file.
	if err := w.file.Close(); err != nil {
		return err
	}

	// Reopen the file in read-write mode and truncate it to clear its content.
	var err error
	w.file, err = os.OpenFile("log.wal", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	err = w.file.Truncate(0)
	if err != nil {
		return err
	}

	// Close and reopen the file for appending.
	if err := w.file.Close(); err != nil {
		return err
	}
	w.file, err = os.OpenFile("log.wal", os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	return nil
}
