package main

import (
	"errors"
)

type MemTable struct {
	table *TreeNode
	size  int
}

func (mem *MemTable) Set(key, value string) error {
	mem.table.insert(Pair{key, value}, true)
	mem.size += len(key) + len(value)
	return nil
}

func (mem *MemTable) Get(key string) (string, error) {
	t := mem.table.search(key)
	if t == nil {
		return "", errors.New("Key probably in the Database")
	}
	if t.marker {
		return t.elem.value, nil
	}
	return "", errors.New("No Such Key in the Database")
}

func (mem *MemTable) Del(key string) (string, error) {
	t := mem.table.search(key)
	if t == nil {
		return "", errors.New("Key probably in the Database")
	}
	if t.marker {
		t.marker = false
		return t.elem.value, nil
	}
	return "", errors.New("No Such Key in the Database")
}

func NewMemTable() *MemTable {
	return &MemTable{
		table: nil,
		size:  0,
	}
}
