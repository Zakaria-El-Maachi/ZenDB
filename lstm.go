package main

import "os"

type Lstm struct {
	mem    *MemTable
	buffer []*MemTable
	wal    *Wal
	// files  map[string]bool
}

func (lstm *Lstm) Set(key, value string) error {
	return lstm.mem.Set(key, value)
}

func (lstm *Lstm) Get(key string) (string, error) {
	v, err := lstm.mem.Get(key)
	if err.Error() == "Key probably in the Database" {
		return "Database Value", nil
	}
	return v, err
}

func (lstm *Lstm) Del(key string) (string, error) {
	return lstm.mem.Del(key)
}

func NewLstm() *Lstm {
	file, err := os.OpenFile("log.wal", os.O_APPEND|os.O_CREATE, 466)
	if err != nil {
		return nil
	}
	var watermark int64 = 0
	return &Lstm{
		mem:    NewMemTable(),
		buffer: make([]*MemTable, 10),
		wal:    &Wal{watermark, file},
	}
}
