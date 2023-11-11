package main

import (
	"log"
	"net/http"
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
	mem  *MemTable
}

func (s Server) fullAddress() string {
	return s.addr + ":" + s.port
}

func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	for k, v := range queries {
		if err := s.mem.Set(k, v[len(v)-1]); err != nil {
			response.WriteHeader(http.StatusNotAcceptable)
			response.Write([]byte(err.Error()))
			return
		}
	}
	response.WriteHeader(http.StatusOK)
}

func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()
	if _, ok := queries[Key]; !ok {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("No keys specified"))
		return
	}
	response.WriteHeader(http.StatusOK)
	for _, k := range queries[Key] {
		if v, err := s.mem.Get(k); err != nil {
			response.Write([]byte(k + " : " + err.Error() + "\n"))
		} else {
			response.Write([]byte(k + " : " + v + "\n"))
		}
	}
}

func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("You are smart Zakaria"))
}

func NewServer() Server {
	s := Server{
		addr: "",
		port: "8081",
		mem:  NewMemTable(),
	}
	http.HandleFunc(Set, s.handleSet)
	http.HandleFunc(Get, s.handleGet)
	http.HandleFunc(Del, s.handleDel)
	log.Fatal(http.ListenAndServe(s.fullAddress(), nil))
	return s
}
