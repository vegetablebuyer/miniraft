package raft

import "io"

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"reflect"
//)

// Command represents an action to be taken on the replicated state machine
type Command interface {
	CommandName() string
}

// CommandApply represents the interface to apply a command to the server
type CommandApply interface {
	Apply(Context) (interface{}, error)
}


type CommandEncoder interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
}


// Create a new instance of a command by name
func newCommand(name string, data []byte) (Command, error) {

}