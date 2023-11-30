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

func decodeHeader(file io.ReadWriteSeeker) (string, uint32, *BloomFilter, uint16, error) {
	file.Seek(0, io.SeekStart)
	p2 := make([]byte, 2)
	p4 := make([]byte, 4)
	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, nil, 0, ErrFileNotEncodedProperly
	}
	magic := string(p4)

	if n, err := file.Read(p4); n != 4 || err != nil {
		return "", 0, nil, 0, ErrFileNotEncodedProperly
	}
	entryCount := binary.LittleEndian.Uint32(p4)

	bitset := make([]byte, 29)
	if n, err := file.Read(bitset); n != 29 || err != nil {
		return "", 0, nil, 0, ErrFileNotEncodedProperly
	}

	file.Read(p2)
	return magic, entryCount, CreateBloomFilter(bitset), binary.LittleEndian.Uint16(p2), nil
}

func parseBody(file io.ReadWriteSeeker, entrycount int) (*MemTable, error) {
	mem := NewMemTable()
	mark := make([]byte, 1)
	var key, value string
	var err error
	h := sha256.New()
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
	magic, entryCount, bloom, _, err := decodeHeader(file)
	if err != nil {
		return "", err
	}
	if magic != MAGIC {
		return "", ErrFileNotRecognized
	}
	if !bloom.Test([]byte(key)) {
		return "", ErrKeyCannotBeInFile
	}
	mem, err := parseBody(file, int(entryCount))
	if err != nil {
		return "", err
	}
	value, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func parse(file io.ReadWriteSeeker, mem *MemTable) error {
	magic, entryCount, _, _, err := decodeHeader(file)
	if err != nil {
		return err
	}
	if magic != MAGIC {
		return ErrFileNotRecognized
	}
	mark := make([]byte, 1)
	var key, value string
	h := sha256.New()
	for i := 0; i < int(entryCount); i++ {
		if _, err = file.Read(mark); err != nil {
			return ErrFileNotEncodedProperly
		}
		key, err = decodeBytes(file)
		if err != nil {
			return err
		}
		if mark[0] == 's' {
			value, err = decodeBytes(file)
			if err != nil {
				return err
			}
			mem.Set(key, value)
			h.Write([]byte(key + value))
		} else if mark[0] == 'd' {
			mem.Del(key)
			h.Write([]byte(key))
		} else {
			return ErrFileNotEncodedProperly
		}

	}

	p := make([]byte, 32)
	if n, err := file.Read(p); n < 32 || err != nil {
		return ErrFileNotEncodedProperly
	}

	if bytes.Compare(h.Sum(nil), p) != 0 {
		return ErrCorruptFile
	}
	return nil
}
