package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	switch pattern {
	case GetPath, DelPath:
		if _, ok := (*queries)[Key]; !ok {
			return errors.New("No specified element to " + pattern[1:])
		}
		if len((*queries)[Key]) != 1 {
			return errors.New("Too many keys to " + pattern[1:] + ", request cancelled")
		}
		return nil
	default:
		return errors.New("Unrecognized pattern")
	}
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
		return errors.New("You should specify one key-value pair to set")
	}
	for key, value := range data {
		if key == "" || !isASCII(key) {
			return fmt.Errorf("Invalid key: %s", key)
		}

		if value == "" || !isASCII(value) {
			return fmt.Errorf("Invalid value for key %s: %s", key, value)
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

func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if err := validate(&queries, GetPath); err != nil {
		writeResponse(&response, StatusBadRequest, err.Error())
		return
	}
	k := queries[Key][0]
	if v, err := s.lstm.Get(k); err != nil {
		writeResponse(&response, StatusBadRequest, k+" : "+err.Error())
	} else {
		writeResponse(&response, StatusOK, k+" : "+v)
	}
}

func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if err := validate(&queries, DelPath); err != nil {
		writeResponse(&response, StatusBadRequest, err.Error())
		return
	}
	k := queries[Key][0]
	if v, err := s.lstm.Del(k); err != nil {
		writeResponse(&response, StatusBadRequest, k+" : "+err.Error())
	} else {
		writeResponse(&response, StatusOK, "Deleted Successfully : "+k+" : "+v)
	}
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
