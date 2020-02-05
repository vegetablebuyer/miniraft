package raftd

import (
	"../cobbler"
	"../raft"
)

// This command writes a value to a key.
type CobblerCommand struct {
	SerialNumber string                   `json:"serialNumber"`
	Args         []map[string]interface{} `json:"args"`
}

// Creates a new write command.
func NewCobblerCommand(serialNumber string, args []map[string]interface{}) *CobblerCommand {
	return &CobblerCommand{
		SerialNumber: serialNumber,
		Args:         args,
	}
}

// The name of the command in the log.
func (c *CobblerCommand) CommandName() string {
	return "cobblerAdd"
}

// Writes a value to a key.
func (c *CobblerCommand) Apply(server raft.Server) (interface{}, error) {
	//db := server.Context().(*DB)
	//db.Put(c.Key, c.Value)
	cli := cobbler.NewClient()
	for _, args := range c.Args {
		cli.EditSystem(c.SerialNumber, args)
	}
	return nil, nil
}
