package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	flushThreshold      = 20
	cleanThreshold      = 50
	CompactionThreshold = 2
	path1               = "Zen_SST\\ZenFile"
	path2               = ".sst"
)

var (
	ErrFileNotRecognized      = errors.New("File not recognized")
	ErrFileNotEncodedProperly = errors.New("File not encoded properly")
	ErrCorruptFile            = errors.New("The file is corrupt")
	ErrKeyNotFound            = errors.New("Key not found")
	ErrKeyDeleted             = errors.New("Key does not exist")
	ErrKeyCannotBeInFile      = errors.New("Key cannot be in current file")
	ErrDeletion               = errors.New("Error While Deleting")
)

// Lstm represents the main storage manager.
type Lstm struct {
	mem      *MemTable
	wal      *Wal
	sstFiles []int
	mu       sync.RWMutex
}

// Set adds a new key-value pair to the storage manager.
func (lstm *Lstm) Set(key, value string) error {
	lstm.mu.Lock()
	defer lstm.mu.Unlock()
	defer lstm.memFlush()
	if err := lstm.wal.RecordSet(key, value); err != nil {
		return err
	}
	return lstm.mem.Set(key, value)
}

func (lstm *Lstm) Search(key string) (string, error) {
	v, err := lstm.mem.Get(key)
	if err != nil && errors.Is(err, ErrKeyNotFound) {
		for i := len(lstm.sstFiles) - 1; i >= 0 && lstm.sstFiles[i] != 0; i-- {
			file, err := os.Open(path1 + fmt.Sprint(lstm.sstFiles[i]) + path2)
			if err != nil {
				log.Println(err)
				continue
			}
			v, err := Search(key, file)
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

// Get retrieves the value associated with a key from the storage manager.
func (lstm *Lstm) Get(key string) (string, error) {
	lstm.mu.RLock()
	defer lstm.mu.RUnlock()
	return lstm.Search(key)
}

// Del removes a key from the storage manager.
func (lstm *Lstm) Del(key string) (string, error) {
	lstm.mu.Lock()
	defer lstm.mu.Unlock()
	defer lstm.memFlush()
	v, err := lstm.Search(key)
	if err == nil {
		if err := lstm.wal.RecordDel(key); err != nil {
			return "", err
		}
		return v, lstm.mem.Del(key)
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

		if err := lstm.wal.Clean(); err != nil {
			log.Println(err)
		}
	}
}

func Recover() (*MemTable, error) {
	log.Println("Recovering...")
	defer log.Println("Recovering Complete\nReady For Requests")
	file, err := os.OpenFile("log.wal", FileFlags, FilePermission)
	if err != nil {
		return nil, err
	}

	mem := NewMemTable()

	mark := make([]byte, 1)
	var key, value string
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	for {
		_, err := file.Read(mark)
		if err == io.EOF {
			break
		}
		if err != nil {
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
		} else if mark[0] == 'd' {
			mem.Del(key)
		} else {
			return nil, ErrFileNotEncodedProperly
		}
	}
	file.Close()
	return mem, nil
}

// LstmDB initializes the storage manager.
func LstmDB() (*Lstm, error) {
	mem, err := Recover()
	if err != nil {
		return nil, err
	}
	sstFiles := getSstFiles()
	file, err := os.OpenFile("log.wal", FileFlags, FilePermission)
	if err != nil {
		return nil, err
	}
	resLstm := &Lstm{
		mem:      mem,
		wal:      &Wal{file},
		sstFiles: sstFiles,
	}
	go resLstm.Compact()
	return resLstm, nil
}

func getSstFiles() []int {
	directory := "Zen_SST"

	var sstFiles []int

	// Read the directory entries
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal("Error reading directory Zen_SST:", err)
	}

	// Define a regular expression to match file names like "ZenFileX.sst"
	re := regexp.MustCompile(`^ZenFile(\d+)\.sst$`)

	// Iterate through the files in the directory
	for _, file := range files {
		match := re.FindStringSubmatch(file.Name())
		if match != nil {
			x, err := strconv.Atoi(match[1])
			if err == nil {
				sstFiles = append(sstFiles, x)
			}
		} else {
			os.Remove(file.Name())
		}
	}
	if len(sstFiles) == 0 {
		sstFiles = append(sstFiles, 0)
	}
	return sstFiles
}

func (lstm *Lstm) Compact() {
	for {
		length := len(lstm.sstFiles)
		if length == 0 {
			continue
		}
		n1 := 0
		n2 := 1
		if lstm.sstFiles[0] == 0 {
			length--
			n1++
			n2++
		}
		if length >= CompactionThreshold {
			lstm.mu.Lock()
			file1, err := os.Open(path1 + fmt.Sprint(lstm.sstFiles[n1]) + path2)
			if err != nil {
				log.Println(err)
			}
			file2, err := os.Open(path1 + fmt.Sprint(lstm.sstFiles[n2]) + path2)
			if err != nil {
				log.Println(err)
			}
			memTemp := NewMemTable()
			if err = Parse(file1, memTemp); err != nil {
				log.Println(err)
			}
			if err = Parse(file2, memTemp); err != nil {
				log.Println(err)
			}
			file1.Close()
			file2.Close()
			os.Remove(file2.Name())
			if err = memTemp.Flush(file1.Name()); err != nil {
				log.Println(err)
			}
			lstm.sstFiles = append(lstm.sstFiles[:n2], lstm.sstFiles[n2+1:]...)
			lstm.mu.Unlock()
		}
	}
}
