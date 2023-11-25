package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const flushThreshhold = 20
const PATH1 = "Zen_SST\\ZenFile"
const PATH2 = ".sst"

type Lstm struct {
	mem      *MemTable
	buffer   []*MemTable
	wal      *Wal
	sstFiles []int
}

func (lstm *Lstm) Set(key, value string) error {
	buffer2 := make([]byte, 2)
	lstm.wal.file.Write([]byte("s"))
	binary.LittleEndian.PutUint16(buffer2, uint16(len(key)))
	lstm.wal.file.Write(buffer2)
	lstm.wal.file.Write([]byte(key))
	binary.LittleEndian.PutUint16(buffer2, uint16(len(value)))
	lstm.wal.file.Write(buffer2)
	lstm.wal.file.Write([]byte(value))
	return lstm.mem.Set(key, value)
}

func (lstm *Lstm) Get(key string) (string, error) {
	v, err := lstm.mem.Get(key)
	if err != nil && err.Error() == "Key probably in the Database" {
		for i := len(lstm.sstFiles) - 1; i > 0; i-- {
			file, err := os.Open(PATH1 + fmt.Sprint(lstm.sstFiles[i]) + PATH2)
			if err != nil {
				log.Println(err)
				continue
			}
			v, err := search(key, file)
			if err != nil {
				if err.Error() == "File not recognized" || err.Error() == "File Not Encoded Properly" || err.Error() == "The File is Corrupt" {
					log.Println(err)
				}
				if err.Error() == "No Such Key in the Database" {
					break
				}
				continue
			}
			return v, nil
		}
		return "", errors.New("No Such Key in the Database")

	}
	return v, err
}

func (lstm *Lstm) Del(key string) (string, error) {
	v, err := lstm.Get(key)
	if err == nil {
		if err0 := lstm.mem.Del(key); err0 != nil {
			return v, errors.New("Error While Deleting")
		}
		buffer2 := make([]byte, 2)
		lstm.wal.file.Write([]byte("d"))
		binary.LittleEndian.PutUint16(buffer2, uint16(len(key)))
		lstm.wal.file.Write(buffer2)
		lstm.wal.file.Write([]byte(key))
	}
	return v, err
}

func (lstm *Lstm) memFlush() {
	for {
		if lstm.mem.size >= flushThreshhold {
			lstm.mem.Flush(PATH1 + fmt.Sprint(lstm.sstFiles[len(lstm.sstFiles)-1]+1) + PATH2)
			lstm.mem = NewMemTable()
			lstm.sstFiles = append(lstm.sstFiles, lstm.sstFiles[len(lstm.sstFiles)-1]+1)
			file, err := os.OpenFile("meta.wal", os.O_APPEND|os.O_CREATE, 466)
			if err != nil {
				log.Fatal(err)
			}
			buffer4 := make([]byte, 4)
			binary.LittleEndian.PutUint32(buffer4, uint32(lstm.sstFiles[len(lstm.sstFiles)-1]))
			file.Write(buffer4)
			lstm.wal.water.Seek(0, io.SeekStart)
			info, err := lstm.wal.file.Stat()
			binary.LittleEndian.PutUint32(buffer4, uint32(info.Size()))
			lstm.wal.water.Write(buffer4)
		}
	}
}

func LstmDB() (*Lstm, error) {
	file, err := os.OpenFile("log.wal", os.O_RDONLY|os.O_CREATE, 466)
	if err != nil {
		return nil, err
	}
	buffer8 := make([]byte, 8)
	buffer4 := make([]byte, 4)
	var exists bool = true
	if _, err = os.Stat("meta.wal"); os.IsNotExist(err) {
		exists = false
	}
	meta, err := os.OpenFile("meta.wal", os.O_RDWR|os.O_CREATE, 466)
	if !exists {
		binary.LittleEndian.PutUint32(buffer4, 0)
		meta.Write(buffer4)
	}
	mem := NewMemTable()
	var watermark int64
	sstFiles := make([]int, 1)
	sstFiles[0] = 0
	exists = true
	if _, err = os.Stat("water.wal"); os.IsNotExist(err) {
		exists = false
	}
	water, err := os.OpenFile("water.wal", os.O_RDWR|os.O_CREATE, 466)
	if !exists {
		watermark = 0
		binary.LittleEndian.PutUint64(buffer8, 0)
		water.Write(buffer8)
	} else {
		water.Read(buffer8)
		watermark = int64(binary.LittleEndian.Uint64(buffer8))
	}

	meta.Seek(4, io.SeekStart)
	for {
		if _, err := meta.Read(buffer4); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		sstFiles = append(sstFiles, int(binary.LittleEndian.Uint32(buffer4)))
	}
	mark := make([]byte, 1)
	var key, value string
	if _, err = file.Seek(watermark, io.SeekStart); err != nil {
		return nil, err
	}
	for {
		_, err = file.Read(mark)
		if err == io.EOF {
			break
		}
		if err != nil {
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
		} else if mark[0] == 'd' {
			mem.Del(key)
		} else {
			return nil, errors.New("File Not Encoded Properly")
		}
	}
	file.Close()
	file, err = os.OpenFile("log.wal", os.O_APPEND, 466)
	if err != nil {
		return nil, err
	}
	resLstm := &Lstm{
		mem:      mem,
		buffer:   make([]*MemTable, 10),
		wal:      &Wal{watermark, file, water, meta},
		sstFiles: sstFiles,
	}
	go resLstm.memFlush()
	return resLstm, nil
}
