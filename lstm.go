package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const flushThreshold = 1024
const cleanThreshold = 1024
const path1 = "Zen_SST\\ZenFile"
const path2 = ".sst"

var (
	ErrFileNotRecognized      = errors.New("File not recognized")
	ErrFileNotEncodedProperly = errors.New("File not encoded properly")
	ErrCorruptFile            = errors.New("The file is corrupt")
	ErrKeyNotFound            = errors.New("Key not found")
	ErrKeyDeleted             = errors.New("Key does not exist")
	ErrKeyCannotBeInFile      = errors.New("Key cannot be in current file")
)

// Lstm represents the main storage manager.
type Lstm struct {
	mem      *MemTable
	wal      *Wal
	sstFiles []int
	mu       sync.Mutex
}

// Set adds a new key-value pair to the storage manager.
func (lstm *Lstm) Set(key, value string) error {
	defer lstm.wal.Clean()
	defer lstm.memFlush()
	buffer2 := make([]byte, 2)
	if _, err := lstm.wal.file.Write([]byte("s")); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(buffer2, uint16(len(key)))
	if _, err := lstm.wal.file.Write(buffer2); err != nil {
		return err
	}
	if _, err := lstm.wal.file.Write([]byte(key)); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(buffer2, uint16(len(value)))
	if _, err := lstm.wal.file.Write(buffer2); err != nil {
		return err
	}
	if _, err := lstm.wal.file.Write([]byte(value)); err != nil {
		return err
	}
	return lstm.mem.Set(key, value)
}

// Get retrieves the value associated with a key from the storage manager.
func (lstm *Lstm) Get(key string) (string, error) {
	v, err := lstm.mem.Get(key)
	if err != nil && errors.Is(err, ErrKeyNotFound) {
		for i := len(lstm.sstFiles) - 1; i > 0; i-- {
			file, err := os.Open(path1 + fmt.Sprint(lstm.sstFiles[i]) + path2)
			if err != nil {
				log.Println(err)
				continue
			}
			v, err := search(key, file)
			if err != nil {
				if errors.Is(err, ErrFileNotRecognized) || errors.Is(err, ErrFileNotEncodedProperly) || errors.Is(err, ErrCorruptFile) {
					log.Println(err)
				}
				if errors.Is(err, ErrKeyDeleted) {
					break
				}
				continue
			}
			return v, nil
		}
		return "", ErrKeyDeleted
	}
	return v, err
}

// Del removes a key from the storage manager.
func (lstm *Lstm) Del(key string) (string, error) {
	defer lstm.memFlush()
	defer lstm.wal.Clean()
	v, err := lstm.Get(key)
	if err == nil {
		if err := lstm.mem.Del(key); err != nil {
			return v, errors.New("Error While Deleting")
		}
		buffer2 := make([]byte, 2)
		if _, err := lstm.wal.file.Write([]byte("d")); err != nil {
			return v, err
		}
		binary.LittleEndian.PutUint16(buffer2, uint16(len(key)))
		if _, err := lstm.wal.file.Write(buffer2); err != nil {
			return v, err
		}
		if _, err := lstm.wal.file.Write([]byte(key)); err != nil {
			return v, err
		}
	}
	return v, err
}

// memFlush periodically flushes the in-memory table to disk. Also handles Wal Cleaning.
func (lstm *Lstm) memFlush() {
	if lstm.mem.size >= flushThreshold {
		if err := lstm.mem.Flush(path1 + fmt.Sprint(lstm.sstFiles[len(lstm.sstFiles)-1]+1) + path2); err != nil {
			log.Println(err)
		}
		lstm.mem = NewMemTable()
		lstm.sstFiles = append(lstm.sstFiles, lstm.sstFiles[len(lstm.sstFiles)-1]+1)

		buffer4 := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer4, uint32(lstm.sstFiles[len(lstm.sstFiles)-1]))
		if _, err := lstm.wal.meta.Write(buffer4); err != nil {
			log.Println(err)
		}

		buffer8 := make([]byte, 8)
		lstm.wal.water.Seek(0, io.SeekStart)
		info, err := lstm.wal.file.Stat()
		if err != nil {
			log.Println(err)
		}
		binary.LittleEndian.PutUint64(buffer8, uint64(info.Size()))
		if _, err := lstm.wal.water.Write(buffer8); err != nil {
			log.Println(err)
		}
	}
}

func (lstm *Lstm) walClean() {
	stats, err := lstm.wal.file.Stat()
	if err != nil {
		return
	}

	if stats.Size() >= cleanThreshold {
		err = lstm.wal.Clean()
		if err != nil {
			log.Printf("Error Cleaning the Wal: %v", err)
		}
	}

	lstm.wal.file, err = os.OpenFile("log.wal", FileFlags, FilePermission)
}

// LstmDB initializes the storage manager.
func LstmDB() (*Lstm, error) {
	file, err := os.OpenFile("log.wal", os.O_RDONLY|os.O_CREATE, FilePermission)
	if err != nil {
		return nil, err
	}
	buffer8 := make([]byte, 8)
	buffer4 := make([]byte, 4)
	var exists bool = true
	if _, err = os.Stat("meta.wal"); os.IsNotExist(err) {
		exists = false
	}
	meta, err := os.OpenFile("meta.wal", FileFlags, FilePermission)
	if !exists {
		binary.LittleEndian.PutUint32(buffer4, 0)
		if _, err := meta.Write(buffer4); err != nil {
			log.Println(err)
		}
	}
	mem := NewMemTable()
	var watermark int64
	sstFiles := make([]int, 1)
	sstFiles[0] = 0
	exists = true
	if _, err = os.Stat("water.wal"); os.IsNotExist(err) {
		exists = false
	}
	water, err := os.OpenFile("water.wal", os.O_RDWR|os.O_CREATE, FilePermission)
	if !exists {
		watermark = 0
		binary.LittleEndian.PutUint64(buffer8, 0)
		if _, err := water.Write(buffer8); err != nil {
			log.Println(err)
		}
	} else {
		if _, err := water.Read(buffer8); err != nil {
			log.Println(err)
		}
		watermark = int64(binary.LittleEndian.Uint64(buffer8))
	}

	if _, err := meta.Seek(4, io.SeekStart); err != nil {
		log.Println(err)
	}
	for {
		if _, err := meta.Read(buffer4); err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return nil, err
		}
		sstFiles = append(sstFiles, int(binary.LittleEndian.Uint32(buffer4)))
	}
	mark := make([]byte, 1)
	var key, value string
	if _, err := file.Seek(watermark, io.SeekStart); err != nil {
		log.Println(err)
		return nil, err
	}
	for {
		_, err := file.Read(mark)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return nil, ErrFileNotEncodedProperly
		}
		key, err = decodeBytes(file)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if mark[0] == 's' {
			value, err = decodeBytes(file)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			mem.Set(key, value)
		} else if mark[0] == 'd' {
			mem.Del(key)
		} else {
			log.Println(err)
			return nil, ErrFileNotEncodedProperly
		}
	}
	file.Close()
	file, err = os.OpenFile("log.wal", FileFlags, FilePermission)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resLstm := &Lstm{
		mem:      mem,
		wal:      &Wal{watermark: watermark, file: file, water: water, meta: meta},
		sstFiles: sstFiles,
	}
	return resLstm, nil
}
