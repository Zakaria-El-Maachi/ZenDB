package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestLstmSetGet tests the Set and Get methods of Lstm.
func TestLstmSetGet(t *testing.T) {
	os.Remove("log.wal")
	os.Remove("meta.wal")
	os.Remove("water.wal")
	lstm, err := LstmDB()
	if err != nil {
		t.Fatalf("Error creating Lstm: %v", err)
	}

	key := "testKey"
	value := "testValue"

	err = lstm.Set(key, value)
	if err != nil {
		t.Errorf("Error setting key-value pair: %v", err)
	}

	result, err := lstm.Get(key)
	if err != nil {
		t.Errorf("Error getting value for key: %v", err)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}

	// Clean up
	defer os.Remove("log.wal")
	defer os.Remove("meta.wal")
	defer os.Remove("water.wal")
	for _, file := range lstm.sstFiles {
		os.Remove(path1 + fmt.Sprint(file) + path2)
	}
}

// TestLstmDel tests the Del method of Lstm.
func TestLstmDel(t *testing.T) {
	os.Remove("log.wal")
	os.Remove("meta.wal")
	os.Remove("water.wal")
	lstm, err := LstmDB()
	if err != nil {
		t.Fatalf("Error creating Lstm: %v", err)
	}

	key := "testKey"
	value := "testValue"

	err = lstm.Set(key, value)
	if err != nil {
		t.Errorf("Error setting key-value pair: %v", err)
	}

	v, err := lstm.Del(key)
	if err != nil {
		t.Errorf("Error deleting key: %v", err)
	}

	if v != value {
		t.Errorf("Expected deleted value %s, got %s", value, v)
	}

	// Clean up
	defer os.Remove("log.wal")
	defer os.Remove("meta.wal")
	defer os.Remove("water.wal")
	for _, file := range lstm.sstFiles {
		os.Remove(path1 + fmt.Sprint(file) + path2)
	}
}

// TestLstmMemFlush tests the memFlush method of Lstm.
func TestLstmMemFlush(t *testing.T) {
	os.Remove("log.wal")
	os.Remove("meta.wal")
	os.Remove("water.wal")
	lstm, err := LstmDB()
	if err != nil {
		t.Fatalf("Error creating Lstm: %v", err)
	}

	key := "testKey"
	value := "testValue"

	for i := 0; i < 1100; i++ {
		err = lstm.Set(key+fmt.Sprint(i), value)
		if err != nil {
			t.Errorf("Error setting key-value pair: %v", err)
		}
	}

	// Allow time for memFlush to execute
	time.Sleep(2 * time.Second)

	// Check if the SST file is created
	_, err = os.Stat(path1 + "1" + path2)
	if err != nil {
		t.Errorf("Error checking SST file: %v", err)
	}

	// Clean up
	defer os.Remove("log.wal")
	defer os.Remove("meta.wal")
	defer os.Remove("water.wal")
	for _, file := range lstm.sstFiles {
		os.Remove(path1 + fmt.Sprint(file) + path2)
	}
}

// TestLstmGetAfterFlush tests the Get method of Lstm after memFlush.
func TestLstmGetAfterFlush(t *testing.T) {
	os.Remove("log.wal")
	os.Remove("meta.wal")
	os.Remove("water.wal")
	lstm, err := LstmDB()
	if err != nil {
		t.Fatalf("Error creating Lstm: %v", err)
	}

	key := "testKey"
	value := "testValue"

	err = lstm.Set(key, value)
	if err != nil {
		t.Errorf("Error setting key-value pair: %v", err)
	}

	// Allow time for memFlush to execute
	time.Sleep(2 * time.Second)

	result, err := lstm.Get(key)
	if err != nil {
		t.Errorf("Error getting value for key after flush: %v", err)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}

	// Clean up
	defer os.Remove("log.wal")
	defer os.Remove("meta.wal")
	defer os.Remove("water.wal")
	for _, file := range lstm.sstFiles {
		os.Remove(path1 + fmt.Sprint(file) + path2)
	}
}
