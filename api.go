package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"unicode"
)

const (
	SetPath = "/set"
	GetPath = "/get"
	DelPath = "/del"
	Key     = "key"
)

const (
	StatusOK               = http.StatusOK
	StatusMethodNotAllowed = http.StatusMethodNotAllowed
	StatusBadRequest       = http.StatusBadRequest
)

var (
	ErrSpecifyToGet = errors.New("No specified key to get")
	ErrSpecifyToDel = errors.New("No specified key to del")
	ErrSpecifyToSet = errors.New("You should specify one key value pair to set")
	ErrTooManyKeys  = errors.New("Too many keys specified, request cancelled")
	ErrInvalidKey   = errors.New("Invalid key")
	ErrInvalidValue = errors.New("Invalid value")
)

type Server struct {
	addr string
	port string
	lstm *Lstm
}

func (s Server) fullAddress() string {
	return s.addr + ":" + s.port
}

func writeResponse(response *http.ResponseWriter, status int, message string) {
	(*response).WriteHeader(status)
	(*response).Write([]byte(message))
}

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

func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

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

func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	helperGetDel(&response, request, s.lstm.Get, "")
}

func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	helperGetDel(&response, request, s.lstm.Del, "Deleted Successfully : ")
}

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
