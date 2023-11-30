package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock Lstm implementation for testing
type mockLstm struct {
	data map[string]string
}

func (m *mockLstm) Set(key, value string) error {
	m.data[key] = value
	return nil
}

func (m *mockLstm) Get(key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", ErrKeyNotFound
}

func (m *mockLstm) Del(key string) (string, error) {
	if val, ok := m.data[key]; ok {
		delete(m.data, key)
		return val, nil
	}
	return "", ErrKeyNotFound
}

func TestHandleSet(t *testing.T) {
	mock := &mockLstm{data: make(map[string]string)}
	server := &Server{lstm: mock}
	inputJSON := `{"testKey": "testValue"}`

	req, err := http.NewRequest("POST", SetPath, strings.NewReader(inputJSON))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleSet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, StatusOK)
	}

	expected := "The key-value pair was set successfully"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	if value, ok := mock.data["testKey"]; !ok || value != "testValue" {
		t.Errorf("Lstm Set method not called correctly")
	}
}

func TestHelperGetDel(t *testing.T) {
	mock := &mockLstm{data: map[string]string{"testKey": "testValue"}}
	server := &Server{lstm: mock}
	req, err := http.NewRequest("GET", GetPath+"?key=testKey", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		helperGetDel(&response, request, server.lstm.Get, "")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, StatusOK)
	}

	expected := "testKey : testValue"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleGet(t *testing.T) {
	mock := &mockLstm{data: map[string]string{"testKey": "testValue"}}
	server := &Server{lstm: mock}
	req, err := http.NewRequest("GET", GetPath+"?key=testKey", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleGet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, StatusOK)
	}

	expected := "testKey : testValue"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleDel(t *testing.T) {
	mock := &mockLstm{data: map[string]string{"testKey": "testValue"}}
	server := &Server{lstm: mock}
	req, err := http.NewRequest("DELETE", DelPath+"?key=testKey", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleDel)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, StatusOK)
	}

	// expected := "Deleted Successfully : testKey"
	// if rr.Body.String() != expected {
	// 	t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	// }

	if _, ok := mock.data["testKey"]; ok {
		t.Errorf("Lstm Del method not called correctly")
	}
}
