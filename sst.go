package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

const MAGIC = "zeni"

func decodeBytes(file io.ReadWriteSeeker) (string, error) {
	var length uint16
	err := binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		return "", ErrFileNotEncodedProperly
	}

	data := make([]byte, length)
	_, err = io.ReadFull(file, data)
	if err != nil {
		return "", ErrFileNotEncodedProperly
	}
	return string(data), nil
}

func decodeHeader(file io.ReadWriteSeeker) (string, uint32, string, uint16, error) {
	file.Seek(0, io.SeekStart)
	p2 := make([]byte, 2)
	p4 := make([]byte, 4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, ErrFileNotEncodedProperly
	}
	magic := string(p4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, ErrFileNotEncodedProperly
	}
	entryCount := binary.LittleEndian.Uint32(p4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, "", 0, ErrFileNotEncodedProperly
	}
	maxOffset := binary.LittleEndian.Uint32(p4)
	if _, err := file.Seek(int64(maxOffset+1), io.SeekStart); err != nil {
		return "", 0, "", 0, ErrFileNotEncodedProperly
	}
	maxKey, err := decodeBytes(file)
	if err != nil {
		return "", 0, "", 0, err
	}
	file.Seek(12, io.SeekStart)
	file.Read(p2)
	return magic, entryCount, maxKey, binary.LittleEndian.Uint16(p2), nil
}

func parse(file io.ReadWriteSeeker, entrycount int) (*MemTable, error) {
	mem := NewMemTable()
	mark := make([]byte, 1)
	var key, value string
	var err error
	h := sha256.New()
	if _, err = file.Seek(14, io.SeekStart); err != nil {
		return nil, ErrFileNotEncodedProperly
	}
	for i := 0; i < entrycount; i++ {
		if _, err = file.Read(mark); err != nil {
			return nil, ErrFileNotEncodedProperly
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
			return nil, ErrFileNotEncodedProperly
		}

	}

	p := make([]byte, 32)
	if n, err := file.Read(p); n < 32 || err != nil {
		return nil, ErrFileNotEncodedProperly
	}

	if bytes.Compare(h.Sum(nil), p) != 0 {
		return nil, ErrCorruptFile
	}
	return mem, nil
}

func search(key string, file io.ReadWriteSeeker) (string, error) {
	magic, entryCount, maxKey, _, err := decodeHeader(file)
	if err != nil {
		return "", err
	}
	if magic != MAGIC {
		return "", ErrFileNotRecognized
	}
	if key > maxKey {
		return "", ErrKeyCannotBeInFile
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
