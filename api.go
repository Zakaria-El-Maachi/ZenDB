package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
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

func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if err := validate(&queries, Set); err != nil {
		writeResponse(&response, http.StatusBadRequest, err.Error())
		return
	}
	for k, v := range queries {
		if len(v) != 1 {
			writeResponse(&response, http.StatusBadRequest, "You should specify one key-value pair to set")
			return
		}
		if err := s.lstm.Set(k, v[0]); err != nil {
			writeResponse(&response, http.StatusBadRequest, err.Error())
			return
		}
	}
	writeResponse(&response, http.StatusOK, "The key value pair was set successfully")
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
	s := Server{
		addr: "",
		port: "8081",
		lstm: NewLstm(),
	}
	http.HandleFunc(Set, s.handleSet)
	http.HandleFunc(Get, s.handleGet)
	http.HandleFunc(Del, s.handleDel)
	log.Fatal(http.ListenAndServe(s.fullAddress(), nil))
	return s
}
