package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"os"
)

// Constants for file-related operations.
const (
	MAGIC                  = "zeni"
	FileReadOnlyPermission = 0444
	FilePermission         = 0666
	FileFlags              = os.O_APPEND | os.O_CREATE | os.O_RDWR
)

func encodeString(input string) []byte {
	// Encode the length of the string using 4 bytes
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, uint16(len(input)))

	// Encode each character in the string
	characterBytes := []byte(input)

	// Combine the length bytes and character bytes

	return append(lengthBytes, characterBytes...)
}

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

func parseBody(file io.ReadWriteSeeker, entrycount int, mem *MemTable) error {
	mark := make([]byte, 1)
	var key, value string
	var err error
	h := sha256.New()
	for i := 0; i < entrycount; i++ {
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

func Search(key string, file io.ReadWriteSeeker) (string, error) {
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
	mem := NewMemTable()
	err = parseBody(file, int(entryCount), mem)
	if err != nil {
		return "", err
	}
	value, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func Parse(file io.ReadWriteSeeker, mem *MemTable) error {
	magic, entryCount, _, _, err := decodeHeader(file)
	if err != nil {
		return err
	}
	if magic != MAGIC {
		return ErrFileNotRecognized
	}
	return parseBody(file, int(entryCount), mem)
}
