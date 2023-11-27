package main

import (
	"crypto/sha256"
	"encoding/binary"
	"os"
)

// Constants for file-related operations.
const (
	FileReadOnlyPermission = 444
	FilePermission         = 666
	FileFlags              = os.O_APPEND | os.O_CREATE
)

// MemTable represents an in-memory table.
type MemTable struct {
	table *TreeNode
	size  int
}

// Set adds a new key-value pair to the in-memory table.
func (mem *MemTable) Set(key, value string) error {
	mem.size += Insert(&mem.table, Pair{marker: true, key: key, value: value})
	return nil
}

// Get retrieves the value associated with a key from the in-memory table.
func (mem *MemTable) Get(key string) (string, error) {
	t := mem.table.Search(key)
	if t == nil {
		return "", ErrKeyNotFound
	}
	if t.elem.marker {
		return t.elem.value, nil
	}
	return "", ErrKeyDeleted
}

// Del removes a key from the in-memory table.
func (mem *MemTable) Del(key string) error {
	p := Pair{marker: false, key: key, value: ""}
	mem.size += Insert(&mem.table, p)
	return nil
}

// Flush writes the contents of the in-memory table to a file.
func (mem *MemTable) Flush(fileName string) error {
	kv := mem.table.Traverse()
	maxOff := mem.table.GetMaxOffset(kv)
	file, err := os.OpenFile(fileName, FileFlags, FilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing to the file with error handling.
	if err := writeToFile(file, MAGIC); err != nil {
		return err
	}
	if err := writeUint32ToFile(file, uint32(mem.table.size)); err != nil {
		return err
	}
	if err := writeUint32ToFile(file, uint32(maxOff)); err != nil {
		return err
	}
	if err := writeUint16ToFile(file, uint16(1)); err != nil {
		return err
	}

	h := sha256.New()
	for _, p := range kv {
		if p.marker {
			if err := writeToFile(file, "s"); err != nil {
				return err
			}
			if err := writeUint16ToFile(file, uint16(len(p.key))); err != nil {
				return err
			}
			if err := writeToFile(file, p.key); err != nil {
				return err
			}
			h.Write([]byte(p.key))
			if err := writeUint16ToFile(file, uint16(len(p.value))); err != nil {
				return err
			}
			if err := writeToFile(file, p.value); err != nil {
				return err
			}
			h.Write([]byte(p.value))
		} else {
			if err := writeToFile(file, "d"); err != nil {
				return err
			}
			if err := writeUint16ToFile(file, uint16(len(p.key))); err != nil {
				return err
			}
			if err := writeToFile(file, p.key); err != nil {
				return err
			}
			h.Write([]byte(p.key))
		}
	}
	if _, err := file.Write(h.Sum(nil)); err != nil {
		return err
	}
	return nil
}

// NewMemTable creates a new in-memory table.
func NewMemTable() *MemTable {
	return &MemTable{
		table: nil,
		size:  0,
	}
}

// Helper function to write a string to a file with error handling.
func writeToFile(file *os.File, data string) error {
	_, err := file.Write([]byte(data))
	return err
}

// Helper function to write a uint32 to a file with error handling.
func writeUint32ToFile(file *os.File, data uint32) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, data)
	_, err := file.Write(buffer)
	return err
}

// Helper function to write a uint16 to a file with error handling.
func writeUint16ToFile(file *os.File, data uint16) error {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, data)
	_, err := file.Write(buffer)
	return err
}
