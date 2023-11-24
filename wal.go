package main

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
)

const BufferSize = 2048

type Wal struct {
	watermark int64
	file      *os.File
	water     *os.File
	meta      *os.File
}

// Writes to the Wal in Append Mode and Flushes immediately
func (w *Wal) Write(op []byte) error {
	if _, err := w.file.Write(op); err != nil {
		return err
	}
	if err := w.file.Sync(); err != nil {
		return err
	}
	return nil
}

// Cleans the WAL from the watermarked entries while updating the watermark
func (w *Wal) Clean() error {
	temp := strings.Split(w.file.Name(), ".")
	file, err := os.Create(temp[0] + "2." + temp[1])
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	b := make([]byte, BufferSize)
	w.file.Seek(w.watermark, io.SeekStart)
	for {
		n, err := w.file.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n, err = file.Write(b[:n]); err != nil {
			return err
		}
	}
	if _, err := io.Copy(w.file, file); err != nil {
		return err
	} else {
		w.watermark = 0
		buffer4 := make([]byte, 4)
		binary.LittleEndian.PutUint64(buffer4, 0)
		w.water.Seek(0, io.SeekStart)
		w.water.Write(buffer4)
		return nil
	}
}
