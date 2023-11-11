package main

import (
	"errors"
)

type MemTable struct {
	table   map[string]string
	size    int
	deleted map[string]bool
}

func (mem *MemTable) Set(key, value string) error {
	mem.table[key] = value
	mem.size += len(key) + len(value)
	return nil
}

func (mem *MemTable) Get(key string) (string, error) {
	if _, ok := mem.deleted[key]; ok {
		return "", errors.New("No Such Key in the Database")
	}
	if v, ok := mem.table[key]; ok {
		return v, nil
	}
	return "", errors.New("Key probably in the Database")
}

func (mem *MemTable) Del(key string) (string, error) {
	value, err := mem.Get(key)
	if err != nil {
		return value, err
	}
	delete(mem.table, key)
	mem.deleted[key] = true
	return value, nil
}

func NewMemTable() *MemTable {
	table := make(map[string]string)
	deleted := make(map[string]bool)
	size := 0
	return &MemTable{
		table,
		size,
		deleted,
	}
}
