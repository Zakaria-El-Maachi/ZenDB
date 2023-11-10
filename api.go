package main

import (
	"log"
	"net/http"
)

const (
	Set = "/set"
	Get = "/get"
	Del = "/del"
)

type Server struct {
	addr string
	port string
}

func (s Server) fullAddress() string {
	return s.addr + ":" + s.port
}

func (s *Server) handleSet(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("You are the best Zakaria"))
}

func (s *Server) handleGet(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("You are cute Zakaria"))
}

func (s *Server) handleDel(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("You are smart Zakaria"))
}

func NewServer() Server {
	s := Server{
		addr: "",
		port: "8081",
	}
	http.HandleFunc(Set, s.handleSet)
	http.HandleFunc(Get, s.handleGet)
	http.HandleFunc(Del, s.handleDel)
	log.Fatal(http.ListenAndServe(s.fullAddress(), nil))
	return s
}
