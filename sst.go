package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
)

const MAGIC = "zeni"

func decodeBytes(file io.ReadWriteSeeker) (string, error) {
	p := make([]byte, 2)
	n, err := file.Read(p)
	if err != nil {
		return "", err
	}
	if n != 2 {
		return "", errors.New("File Not Encoded Properly")
	}
	length := binary.LittleEndian.Uint16(p)
	p = make([]byte, length)
	n, err = file.Read(p)
	if err != nil {
		return "", err
	}
	if uint16(n) != length {
		return "", errors.New("File Not Encoded Properly")
	}
	return string(p), nil
}

func decodeHeader(file io.ReadWriteSeeker) (uint32, uint32, string, uint32, error) {
	file.Seek(0, io.SeekStart)
	p2 := make([]byte, 2)
	p4 := make([]byte, 4)
	file.Read(p4)
	magic := binary.LittleEndian.Uint32(p4)
	file.Read(p4)
	entryCount := binary.LittleEndian.Uint32(p4)
	file.Read(p4)
	maxOffset := binary.LittleEndian.Uint32(p4)
	file.Seek(int64(maxOffset), io.SeekStart)
	maxKey, err := decodeBytes(file)
	if err != nil {
		return 0, 0, "", 0, err
	}
	file.Seek(12, io.SeekStart)
	file.Read(p2)
	return magic, entryCount, maxKey, binary.LittleEndian.Uint32(p2), nil
}

func parse(file io.ReadWriteSeeker, entrycount int) (*MemTable, error) {
	mem := NewMemTable()
	h := sha256.New()
	file.Seek(0, 14)
	var key, value string
	var err error
	for i := 0; i < entrycount; i++ {
		key, err = decodeBytes(file)
		if err != nil {
			return nil, err
		}
		value, err = decodeBytes(file)
		if err != nil {
			return nil, err
		}
		if err = mem.Set(key, value); err != nil {
			return nil, err
		}
		h.Write([]byte(key + value))
	}

	p := make([]byte, 32)
	file.Read(p)

	if bytes.Compare(h.Sum(nil), p) != 0 {
		return nil, errors.New("The File is Corrupt")
	}
	return mem, nil
}

func search(key string, file io.ReadWriteSeeker) (bool, error) {
	magic, entryCount, maxKey, version, err := decodeHeader(file)
	if err != nil {
		return false, err
	}
	if key > maxKey {
		return false, nil
	}
	mem := parse(file)
}
