package raft

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	Stopped 		= "stopped"
	Initialized 	= "initialized"
	Follower 		= "Follower"
	Candidate 		= "candidate"
	Leader 			= "leader"
	Snapshotting 	= "snapshotting"
)

const (
	MaxLogEntriesPerRequest		= 2000
)

const (
	// DefaultHeartbeatInterval is the interval that the leader will send
	// AppendEntriesRequests to followers to maintain leadership.
	DefaultHeartbeatInterval 	= 50 * time.Microsecond
	DefaultElectionTimeout		= 150 * time.Microsecond
)

type  Server interface {
	Name() string
	Context() interface{}
	StateMachine() StateMachine
	Leader() string
	State() string
	Path() string

	Start() error
}

type raftServer struct {

	name 		string
	path 		string
	state 		string
	context 	interface{}
	currentTerm uint64

	log 	*Log
	leader 	string
	peers 	map[string]*Peer
	mutex 	sync.RWMutex

	stopped chan bool

	electionTimeout 	time.Duration
	heartbeatInterval 	time.Duration

	stateMachine 			StateMachine
	maxLogEntriesPerRequest uint64
	connectionString 		string
}


// An interval event to be processed by the server's event loop
type ev struct {
	target 		interface{}
	returnValue interface{}
	c			chan error
}

func NewRaftServer(name string, path string, stateMachine StateMachine, ctx interface{},
						connectionString string) (Server, error) {
	if name == "" {
		return nil, errors.New("raft Server name cannot br blank")
	}

	s := &raftServer{
		name:						name,
		path:						path,
		stateMachine:				stateMachine,
		context:					ctx,
		state:						Stopped,
		peers:						make(map[string]*Peer),
		log:						newLog(),
		electionTimeout:			DefaultElectionTimeout,
		heartbeatInterval:			DefaultHeartbeatInterval,

		maxLogEntriesPerRequest: 	MaxLogEntriesPerRequest,

		connectionString:			connectionString,
	}

	// Setup log apply function
	s.log.ApplyFunc = func(e *LogEntry, c Command) (interface{}, error) {

		// Apply command to state machine
		switch c := c.(type) {
		case CommandApply:
			return c.Apply(&context{
				server:			s,
				currentTerm:	s.currentTerm,
				currentIndex:	s.log.intervalCurrentIndex(),
				commitIndex:	s.log.commitIndex,
			})
		default:
			return nil, fmt.Errorf("command does not implement Apply()")
		}
	}

	return s, nil
}

func (s *raftServer) Name() string {
	return s.name
}

func (s *raftServer) Context() interface{} {
	return s.context
}

func (s *raftServer) StateMachine() StateMachine {
	return s.stateMachine
}

func (s *raftServer) Leader() string {
	return s.leader
}

func (s *raftServer) State() string {
	return s.state
}

func (s *raftServer) Path() string {
	return s.path
}

// Start the raft server
// If log entries exist then allow promotion to candidate if no AEs received.
// If no log entries exist then wait for AEs from another node.
// If no log entries exist and a self-join command is issued then
// immediately become leader and commit entry.
func (s *raftServer) Start() error {
	if s.IsRunning() {
		return fmt.Errorf("raft server is already running[%v]", s.state)
	}


}

// Initialize the raft server
// If there is no previous log file under the given path, Init() will create an empty log file.
// Otherwise, Init() will load in the log entries from the log file
func (s *raftServer) Init() error {
	if s.IsRunning() {
		return fmt.Errorf("raft server is already running[%v]", s.state)
	}

	// Server has been initialized or server was stopped after initialized
	// If log has been initialized, we know that the server was stopped after
	// running.
	if s.state == Initialized || s.log.initialized {
		s.state = Initialized
		return nil
	}

	// Create snapshot directory if it does not exist
	err := os.Mkdir(path.Join(s.path, "snapshot"), 0700)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("raft: Initialized error: %s", err)
	}

	// Initialize the log and load it up
	if err := s.log.
}


// Check if the server is currently running
func (s *raftServer) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state != Stopped && s.state != Initialized
}


// Retrieves the log path for the server
func (s *raftServer) LogPath() string {
	return path.Join(s.path, "log")
}
