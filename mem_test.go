package main

// import (
// 	"fmt"
// 	"os"
// 	"testing"
// )

// // TestMemTableSetGet tests the Set and Get methods of MemTable.
// func TestMemTableSetGet(t *testing.T) {
// 	mem := NewMemTable()
// 	for i := 0; i < 100000; i++ {
// 		key := fmt.Sprintf("testKey%d", i)
// 		value := fmt.Sprintf("testValue%d", i)

// 		err := mem.Set(key, value)
// 		if err != nil {
// 			t.Errorf("Error setting key-value pair %d: %v", i, err)
// 		}

// 		result, err := mem.Get(key)
// 		if err != nil {
// 			t.Errorf("Error getting value for key %d: %v", i, err)
// 		}

// 		if result != value {
// 			t.Errorf("Expected value %s, got %s", value, result)
// 		}
// 	}
// }

// // TestMemTableDel tests the Del method of MemTable.
// func TestMemTableDel(t *testing.T) {
// 	mem := NewMemTable()
// 	for i := 0; i < 100000; i++ {
// 		key := fmt.Sprintf("testKey%d", i)
// 		value := fmt.Sprintf("testValue%d", i)

// 		err := mem.Set(key, value)
// 		if err != nil {
// 			t.Errorf("Error setting key-value pair %d: %v", i, err)
// 		}

// 		err = mem.Del(key)
// 		if err != nil {
// 			t.Errorf("Error deleting key %d: %v", i, err)
// 		}

// 		_, err = mem.Get(key)
// 		if err == nil {
// 			t.Errorf("Expected error for deleted key %d, got nil", i)
// 		}
// 	}
// }

// // TestMemTableFlush tests the Flush method of MemTable.
// func TestMemTableFlush(t *testing.T) {
// 	mem := NewMemTable()

// 	key := "testKey"
// 	value := "testValue"

// 	err := mem.Set(key, value)
// 	if err != nil {
// 		t.Errorf("Error setting key-value pair: %v", err)
// 	}

// 	fileName := "testFile.sst"
// 	err = mem.Flush(fileName)
// 	if err != nil {
// 		t.Errorf("Error flushing MemTable to file: %v", err)
// 	}

// 	// Clean up test file
// 	defer os.Remove(fileName)
// }
