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
	Set = "/set"
	Get = "/get"
	Del = "/del"
	Key = "key"
)

type Server struct {
	addr string
	port string
	lstm *Lstm
}

func (s Server) fullAddress() string {
	return s.addr + ":" + s.port
}

func writeResponse(response *http.ResponseWriter, a int, b string) {
	(*response).WriteHeader(a)
	(*response).Write([]byte(b))
}

func validate(queries *url.Values, pattern string) error {
	switch pattern {
	case Set:
		if len(*queries) != 1 {
			return errors.New("You should specify one key-value pair to set")
		}
		return nil
	case Get:
		if _, ok := (*queries)[Key]; !ok {
			return errors.New("No specified element to get")
		}
		if len((*queries)[Key]) != 1 {
			return errors.New("Too many keys to get, request cancelled")
		}
		return nil
	case Del:
		if _, ok := (*queries)[Key]; !ok {
			return errors.New("No specified element to delete")
		}
		if len((*queries)[Key]) != 1 {
			return errors.New("Too many keys to delete, request cancelled")
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

// func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
// 	queries := request.URL.Query()
// 	if err := validate(&queries, Set); err != nil {
// 		writeResponse(&response, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	for k, v := range queries {
// 		if len(v) != 1 {
// 			writeResponse(&response, http.StatusBadRequest, "You should specify one key-value pair to set")
// 			return
// 		}
// 		if err := s.lstm.Set(k, v[0]); err != nil {
// 			writeResponse(&response, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 	}
// 	writeResponse(&response, http.StatusOK, "The key value pair was set successfully")
// }

func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeResponse(&response, http.StatusMethodNotAllowed, "Method not allowed. Only POST requests are allowed.")
		return
	}
	var requestBody map[string]string
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		writeResponse(&response, http.StatusBadRequest, "Error decoding JSON data: "+err.Error())
		return
	}
	if err := validateJSON(requestBody); err != nil {
		writeResponse(&response, http.StatusBadRequest, err.Error())
		return
	}

	for k, v := range requestBody {
		if err := s.lstm.Set(k, v); err != nil {
			writeResponse(&response, http.StatusBadRequest, err.Error())
			return
		}
	}
	writeResponse(&response, http.StatusOK, "The key-value pair was set successfully")
}

func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if err := validate(&queries, Get); err != nil {
		writeResponse(&response, http.StatusBadRequest, err.Error())
		return
	}
	k := queries[Key][0]
	if v, err := s.lstm.Get(k); err != nil {
		writeResponse(&response, http.StatusBadRequest, k+" : "+err.Error())
	} else {
		writeResponse(&response, http.StatusOK, k+" : "+v)
	}
}

func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if err := validate(&queries, Del); err != nil {
		writeResponse(&response, http.StatusBadRequest, err.Error())
		return
	}
	k := queries[Key][0]
	if v, err := s.lstm.Del(k); err != nil {
		writeResponse(&response, http.StatusBadRequest, k+" : "+err.Error())
	} else {
		writeResponse(&response, http.StatusOK, "Deleted Successfully : "+k+" : "+v)
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
	http.HandleFunc(Set, s.handleSet)
	http.HandleFunc(Get, s.handleGet)
	http.HandleFunc(Del, s.handleDel)
	log.Fatal(http.ListenAndServe(s.fullAddress(), nil))
	return s
}
