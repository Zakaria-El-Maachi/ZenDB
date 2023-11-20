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
		return "", errors.New("File Not Encoded Properly")
	}
	if uint16(n) != length {
		return "", errors.New("File Not Encoded Properly")
	}
	return string(p), nil
}

func decodeHeader(file io.ReadWriteSeeker) (string, uint32, string, uint32, error) {
	file.Seek(0, io.SeekStart)
	p2 := make([]byte, 2)
	p4 := make([]byte, 4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, errors.New("File Not Encoded Properly")
	}
	magic := string(p4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, errors.New("File Not Encoded Properly")
	}
	entryCount := binary.LittleEndian.Uint32(p4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, errors.New("File Not Encoded Properly")
	}
	maxOffset := binary.LittleEndian.Uint32(p4)
	if _, err := file.Seek(int64(maxOffset), io.SeekStart); err != nil {
		return "", 0, "", 0, errors.New("File Not Encoded Properly")
	}
	maxKey, err := decodeBytes(file)
	if err != nil {
		return "", 0, "", 0, err
	}
	file.Seek(12, io.SeekStart)
	file.Read(p2)
	return magic, entryCount, maxKey, binary.LittleEndian.Uint32(p2), nil
}

func parse(file io.ReadWriteSeeker, entrycount int) (*MemTable, error) {
	mem := NewMemTable()
	mark := make([]byte, 1)
	var key, value string
	var err error
	h := sha256.New()
	if _, err = file.Seek(0, 14); err != nil {
		return nil, errors.New("File Not Encoded Properly")
	}
	for i := 0; i < entrycount; i++ {
		if _, err = file.Read(mark); err != nil {
			return nil, errors.New("File Not Encoded Properly")
		}
		key, err = decodeBytes(file)
		if err != nil {
			return nil, err
		}
		if mark[0] == 's' {
			value, err = decodeBytes(file)
			if err != nil {
				return nil, err
			}
			mem.Set(key, value)
			h.Write([]byte(key + value))
		} else if mark[0] == 'd' {
			mem.Del(key)
			h.Write([]byte(key))
		} else {
			return nil, errors.New("File Not Encoded Properly")
		}

	}

	p := make([]byte, 32)
	if n, err := file.Read(p); n < 32 || err != nil {
		return nil, errors.New("File Not Encoded Properly")
	}

	if bytes.Compare(h.Sum(nil), p) != 0 {
		return nil, errors.New("The File is Corrupt")
	}
	return mem, nil
}

func search(key string, file io.ReadWriteSeeker) (string, error) {
	magic, entryCount, maxKey, _, err := decodeHeader(file)
	if err != nil {
		return "", err
	}
	if magic != MAGIC {
		return "", errors.New("File not recognized")
	}
	if key > maxKey {
		return "", errors.New("Key cannot be in current file")
	}
	mem, err := parse(file, int(entryCount))
	if err != nil {
		return "", err
	}
	value, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}
