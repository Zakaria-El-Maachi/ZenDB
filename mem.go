package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"os"
)

type MemTable struct {
	table *TreeNode
	size  int
}

func (mem *MemTable) Set(key, value string) error {
	mem.table.insert(Pair{true, key, value})
	mem.size += len(key) + len(value)
	return nil
}

func (mem *MemTable) Get(key string) (string, error) {
	t := mem.table.search(key)
	if t == nil {
		return "", errors.New("Key probably in the Database")
	}
	if t.elem.marker {
		return t.elem.value, nil
	}
	return "", errors.New("No Such Key in the Database")
}

func (mem *MemTable) Del(key string) (string, error) {
	t := mem.table.search(key)
	if t == nil {
		p := Pair{marker: false, key: key, value: ""}
		mem.table.insert(p)
	}
	if t.elem.marker {
		t.elem.marker = false
		return t.elem.value, nil
	}
	return "", errors.New("No Such Key in the Database")
}

func (mem *MemTable) Flush(fileName string) error {
	kv := mem.table.traverse()
	maxOff := mem.table.getMaxOffset(kv)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 444)
	if err != nil {
		return err
	}
	buffer4 := make([]byte, 4)
	buffer2 := make([]byte, 2)
	file.Write([]byte(MAGIC))
	binary.LittleEndian.PutUint32(buffer4, uint32(mem.table.size))
	file.Write(buffer4)
	binary.LittleEndian.PutUint32(buffer4, uint32(maxOff))
	file.Write(buffer4)
	binary.LittleEndian.PutUint16(buffer2, uint16(1))
	file.Write(buffer2)

	h := sha256.New()
	for _, p := range kv {
		if p.marker == true {
			file.Write([]byte("s"))
			binary.LittleEndian.PutUint16(buffer2, uint16(len(p.key)))
			file.Write(buffer2)
			file.Write([]byte(p.key))
			h.Write([]byte(p.key))
			binary.LittleEndian.PutUint16(buffer2, uint16(len(p.value)))
			file.Write(buffer2)
			file.Write([]byte(p.value))
			h.Write([]byte(p.value))
		} else {
			file.Write([]byte("d"))
			binary.LittleEndian.PutUint16(buffer2, uint16(len(p.key)))
			file.Write(buffer2)
			file.Write([]byte(p.key))
		}
	}
	file.Write(h.Sum(nil))
	return nil
}

func NewMemTable() *MemTable {
	return &MemTable{
		table: nil,
		size:  0,
	}
}
