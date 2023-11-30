package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"unicode"
)

// Constants representing API paths
const (
	SetPath = "/set"
	GetPath = "/get"
	DelPath = "/del"
	Key     = "key"
)

// Constants representing HTTP response status codes
const (
	StatusOK               = http.StatusOK
	StatusMethodNotAllowed = http.StatusMethodNotAllowed
	StatusBadRequest       = http.StatusBadRequest
)

// Custom error messages
var (
	ErrSpecifyToGet = errors.New("No specified key to get")
	ErrSpecifyToDel = errors.New("No specified key to del")
	ErrSpecifyToSet = errors.New("You should specify one key value pair to set")
	ErrTooManyKeys  = errors.New("Too many keys specified, request cancelled")
	ErrInvalidKey   = errors.New("Invalid key")
	ErrInvalidValue = errors.New("Invalid value")
)

type DB interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) (string, error)
}

type Server struct {
	addr string
	port string
	lstm DB
}

// fullAddress returns the full address of the server.
func (s Server) fullAddress() string {
	return s.addr + ":" + s.port
}

// writeResponse writes an HTTP response with the given status code and message.
func writeResponse(response *http.ResponseWriter, status int, message string) {
	(*response).WriteHeader(status)
	(*response).Write([]byte(message))
}

// validate checks if the specified key is present in the URL queries and follows a specified pattern.
func validate(queries *url.Values, pattern string) error {
	var err error
	switch pattern {
	case GetPath:
		err = ErrSpecifyToGet
	case DelPath:
		err = ErrSpecifyToDel
	}
	if _, ok := (*queries)[Key]; !ok {
		return err
	}
	if len((*queries)[Key]) != 1 {
		return ErrTooManyKeys
	}
	k := (*queries)[Key][0]
	if k == "" || !isASCII(k) {
		return ErrInvalidKey
	}
	return nil
}

// isASCII checks if a string contains only ASCII characters.
func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// validateJSON checks if the JSON data contains a single key-value pair with valid characters.
func validateJSON(data map[string]string) error {
	if len(data) != 1 {
		return ErrSpecifyToSet
	}
	for key, value := range data {
		if key == "" || !isASCII(key) {
			return ErrInvalidKey
		}

		if value == "" || !isASCII(value) {
			return ErrInvalidValue
		}
	}

	return nil
}

// handleSet handles the "/set" endpoint, setting key-value pairs in the storage.
func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeResponse(&response, StatusMethodNotAllowed, "Method not allowed. Only POST requests are allowed.")
		return
	}
	var requestBody map[string]string
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		writeResponse(&response, StatusBadRequest, "Error decoding JSON data: "+err.Error())
		return
	}
	if err := validateJSON(requestBody); err != nil {
		writeResponse(&response, StatusBadRequest, err.Error())
		return
	}

	for k, v := range requestBody {
		err := s.lstm.Set(k, v)
		if err != nil {
			writeResponse(&response, StatusBadRequest, err.Error())
		}
		break
	}
	writeResponse(&response, StatusOK, "The key-value pair was set successfully")
}

// helperGetDel is a helper function for handling "/get" and "/del" endpoints.
func helperGetDel(response *http.ResponseWriter, request *http.Request, function func(string) (string, error), format string) {
	queries := request.URL.Query()
	if err := validate(&queries, GetPath); err != nil {
		writeResponse(response, StatusBadRequest, err.Error())
		return
	}
	k := queries[Key][0]
	if v, err := function(k); err != nil {
		writeResponse(response, StatusBadRequest, k+" : "+err.Error())
	} else {
		writeResponse(response, StatusOK, format+k+" : "+v)
	}
}

// handleGet handles the "/get" endpoint, retrieving the value for a specified key.
func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	helperGetDel(&response, request, s.lstm.Get, "")
}

// handleDel handles the "/del" endpoint, deleting a specified key from the storage.
func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	helperGetDel(&response, request, s.lstm.Del, "Deleted Successfully : ")
}

// NewServer creates a new instance of the HTTP server.
func NewServer() Server {
	lstm, err := LstmDB()
	if err != nil {
		log.Fatal(err)
	}
	s := Server{
		addr: "",
		port: "8081",
		lstm: lstm,
	}
	http.HandleFunc(SetPath, s.handleSet)
	http.HandleFunc(GetPath, s.handleGet)
	http.HandleFunc(DelPath, s.handleDel)
	log.Fatal(http.ListenAndServe(s.fullAddress(), nil))
	return s
}
