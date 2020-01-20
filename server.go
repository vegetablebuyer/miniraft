package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"

	"./raft"
)

type Server struct {
	name 		string
	host 		string
	port 		int
	path 		string
	router 		*mux.Router
	httpServer 	*http.Server
	raftServer  raft.Server
	db 			*DB
	mutex 		sync.RWMutex
}

// Creates a new server
func NewServer(path string, host string, port int) *Server {
	s := &Server{
		host: 	host,
		port: 	port,
		path: 	path,
		db:	  	NewDB(),
		router: mux.NewRouter(),
	}

	// Read the name for server, if fails, creates one
	if b, err := ioutil.ReadFile(filepath.Join(path, "name")); err == nil {
		s.name = string(b)
	} else {
		// Create a random name for server and write to the configuration
		s.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		if err = ioutil.WriteFile(filepath.Join(path, "name"), []byte(s.name), 0644); err != nil {
			panic(err)
		}
	}
	return s
}

// Start the server
func (s *Server) ListenAndServe(leader string) error {
	var err error

	log.Printf("Initializing Raft Server: %s", s.path)

	s.raftServer, err = raft.NewRaftServer(s.name, s.path, nil, s.db,"")
	if err != nil {
		log.Fatal(err)
	}

	s.httpServer = &http.Server{
		Addr: 		fmt.Sprintf(":%d", s.port),
		Handler: 	s.router,
	}

	//s.router.HandleFunc("/db/{key}", ).Methods("GET")


	return s.httpServer.ListenAndServe()
}

