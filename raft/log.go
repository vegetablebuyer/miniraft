package raft

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"io"
	"log"
	"os"
	"sync"

	"./protobuf"
)


// A log is a collection of log entries that are persisted to durable storage
type Log struct {
	ApplyFunc 	func(*LogEntry, Command) (interface{}, error)
	file 	  	*os.File
	path 		string
	entries 	[]*LogEntry
	commitIndex uint64
	sync		sync.RWMutex
	startIndex 	uint64	// the index before the first entry in the log entries
	startTerm	uint64
	initialized	bool
}


// Create a new log
func newLog() *Log {
	return &Log{
		entries: make([]*LogEntry, 0),
	}
}

// The current index in the log without locking
func (l *Log) intervalCurrentIndex() uint64 {
	if len(l.entries) == 0 {
		return l.startIndex
	}
	return l.entries[len(l.entries)-1].Index()
}

// Open the log file and reads existing entries. The log can remain open and
// continue to append entries to the end of log
func (l *Log) open(path string) error {
	// Read all the entries from the log if one exists.
	var readBytes int64

	var err error

	l.file, err = os.OpenFile(path, os.O_RDWR, 0600)
	l.path = path

	if err != nil {
		// if the log does not exist before we create the log file
		// and set commitIndex to 0
		if os.IsNotExist(err){
			l.file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
			log.Printf("log.open.create %s", path)
			if err == nil {
				l.initialized = true
			}
			return err
		}
		return err
	}

	// Read the log file and decode entries
	for {
		// Instantiate log entry and decode into it
		entry, _ := newLogEntry(l, nil, 0, 0, nil)
		entry.Position, _ = l.file.Seek(0, os.SEEK_CUR)

		n, err := entry.Decode(l.file)
		if err != nil {
			if err == io.EOF {
				log.Printf("open.log.appned: finish")
			} else {
				if err = os.Truncate(path, readBytes); err != nil {
					return fmt.Errorf("raft.Log: Unable to recover: %v", err)
				}
			}
			break
		}
		if entry.Index() > l.startIndex {
			// Append entry
			l.entries = append(l.entries, entry)
			if entry.Index() <= l.commitIndex {
				command, err = n
			}
		}
	}

}