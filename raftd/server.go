package raftd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"../raft"
)

type Server struct {
	name       string
	host       string
	port       int
	path       string
	router     *mux.Router
	httpServer *http.Server
	raftServer raft.Server
	db         *DB
	mutex      sync.RWMutex
}

// Creates a new server
func NewServer(path string, host string, port int) *Server {
	s := &Server{
		host:   host,
		port:   port,
		path:   path,
		db:     NewDB(),
		router: mux.NewRouter(),
	}

	// Read the name for server, if fails, creates one
	if b, err := ioutil.ReadFile(filepath.Join(path, "name")); err == nil {
		log.Println("read name from file")
		s.name = string(b)
	} else {
		// Create a random name for server and write to the configuration
		log.Println("create name from random")
		s.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		if err = ioutil.WriteFile(filepath.Join(path, "name"), []byte(s.name), 0644); err != nil {
			panic(err)
		}
	}
	log.Println("server name : ", s.name)
	return s
}

// Return the connecting string
func (s *Server) connectionString() string {
	return fmt.Sprintf("http://%s:%d", s.host, s.port)
}

// Start the server
func (s *Server) ListenAndServe(leader string) error {
	var err error

	log.Printf("Initializing Raft Server: %s", s.path)

	transporter := raft.NewHTTPTransporter("/raft", 200*time.Millisecond)
	s.raftServer, err = raft.NewServer(s.name, s.path, transporter, nil, s.db, "")
	if err != nil {
		log.Fatal(err)
	}
	transporter.Install(s.raftServer, s)
	err = s.raftServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	if leader != "" {
		// Join to leader if specified.

		log.Println("Attempting to join leader:", leader)

		if !s.raftServer.IsLogEmpty() {
			log.Fatal("Cannot join with an existing log")
		}
		if err := s.Join(leader); err != nil {
			log.Fatal(err)
		}

	} else if s.raftServer.IsLogEmpty() {
		// Initialing the server by joining it self

		log.Println("Initialing a new cluster")

		_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             s.raftServer.Name(),
			ConnectionString: s.connectionString(),
		})
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Println("Recovered from log")
	}

	log.Println("Initializing HTTP server")

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	s.router.HandleFunc("/db/{key}", s.readHandler).Methods("GET")
	s.router.HandleFunc("/db/{key}", s.writeHandler).Methods("POST")
	s.router.HandleFunc("/cobbler", s.cobblerHandler).Methods("POST")
	s.router.HandleFunc("/cblr/svc/op/upload", s.uploadHandler)
	s.router.HandleFunc("/cblr/svc/op/trig/mode/addpool/system/{sn}", s.updatePoolHandler)
	s.router.HandleFunc("/cblr/svc/op/task/mode/init_finished/task/{job_id}", s.updateTaskHandler)
	s.router.HandleFunc("/join", s.joinHandler).Methods("POST")

	log.Println("Listening at:", s.connectionString())

	return s.httpServer.ListenAndServe()
}

// Joins to the leader of an existing cluster.
func (s *Server) Join(leader string) error {
	command := raft.DefaultJoinCommand{
		Name:             s.raftServer.Name(),
		ConnectionString: s.connectionString(),
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(command)
	resp, err := http.Post(fmt.Sprintf("http://%s/join", leader), "application/json", &b)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil

}

// This is a hack around Gorilla mux not providing the correct net/http
// HandleFunc() interface.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

func (s *Server) joinHandler(w http.ResponseWriter, req *http.Request) {
	command := &raft.DefaultJoinCommand{}
	if err := json.NewDecoder(req.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := s.raftServer.Do(command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) readHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	value := s.db.Get(vars["key"])
	_, _ = w.Write([]byte(value))
}

func (s *Server) writeHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value := string(b)

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(NewWriteCommand(vars["key"], value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) updatePoolHandler(w http.ResponseWriter, req *http.Request) {
	leaderName := s.raftServer.Leader()
	if leaderName != s.raftServer.Name() {
		s.redirectToLeader(w, req)
		return
	}
	serialNumber := req.Header.Get("X_SERIALNUM")
	ip := req.Header.Get("X_IP")
	updateCommand := NewUpdateCommand(serialNumber, serialNumber, "pool", "system", ip)
	_, err := s.raftServer.Do(updateCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) updateTaskHandler(w http.ResponseWriter, req *http.Request) {
	leaderName := s.raftServer.Leader()
	if leaderName != s.raftServer.Name() {
		s.redirectToLeader(w, req)
		return
	}
	vars := mux.Vars(req)
	jobId := vars["job_id"]
	serialNumber := req.Header.Get("X_SERIALNUM")
	ip := req.Header.Get("X_IP")

	updateCommand := NewUpdateCommand(serialNumber, jobId, "task", jobId, ip)
	_, err := s.raftServer.Do(updateCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) uploadHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("start upload report file")
	leaderName := s.raftServer.Leader()
	if leaderName != s.raftServer.Name() {
		s.redirectToLeader(w, req)
		return
	}
	serialNumber := req.Header.Get("name")
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	info := string(b)

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(NewUploadCommand(serialNumber, info))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) cobblerHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("processing cobbler request:")

	// If i was not the leader, transmit the request to the leader to process
	leaderName := s.raftServer.Leader()
	if leaderName != s.raftServer.Name() {
		s.redirectToLeader(w, req)
		return
	}

	// Read the value from the POST body.
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("args:", string(b))
	args := &CobblerCommand{}

	if err = json.Unmarshal(b, args); err != nil {
		handleResponse(w, err.Error(), http.StatusBadRequest)
	}
	result := &CobblerResult{
		SerialNumber: args.SerialNumber,
	}

	// Execute the command against the Raft server.
	_, err = s.raftServer.Do(NewCobblerCommand(args.SerialNumber, args.Action, args.Args))
	if err != nil {
		result.IsSucceed = false
		result.Result = err.Error()
	} else {
		result.IsSucceed = true
		result.Result = ""
	}
	log.Println(result)
	if b, err := json.Marshal(result); err != nil {
		handleResponse(w, err.Error(), http.StatusBadRequest)
	} else {
		handleResponse(w, string(b), http.StatusOK)
	}
}

func (s *Server) redirectToLeader(w http.ResponseWriter, req *http.Request) {
	leader := s.raftServer.Peers()[s.raftServer.Leader()]
	url := fmt.Sprintf("%v%v", leader.ConnectionString, req.URL)
	log.Println("redirect to :", url)
	if res, err := httpPost(url, "json", req.Body, req.Header); err != nil {
		handleResponse(w, err.Error(), http.StatusBadRequest)
	} else {
		p := make([]byte, res.ContentLength)
		_, _ = res.Body.Read(p)
		handleResponse(w, string(p), res.StatusCode)
	}
}

func handleResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	log.Println(message)
	_, _ = fmt.Fprintln(w, message)
}
