package raft

import (
	"sync"
	"time"
)

// A peer is a reference to another server involved in the raft cluster
type Peer struct {
	server 				*Server
	Name 				string `json:"name"`
	ConnectionString 	string `json:"connectionString"`
	prevLogIndex		uint64
	stopChan			chan bool
	heartbeatInterval	time.Duration
	lastActivity		time.Time
	sync.RWMutex
}


// Create a new peer
func newPeer(server *Server, name string, connectionString string, heartbeatInterval time.Duration) *Peer {
	return &Peer{
		server: 			server,
		Name: 				name,
		ConnectionString:	connectionString,
		heartbeatInterval:	heartbeatInterval,
	}
}


// Set the heartbeat timeout
func (p *Peer) setHeartbeatInterval(duration time.Duration) {
	p.heartbeatInterval = duration
}
