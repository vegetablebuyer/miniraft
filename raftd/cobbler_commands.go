package raftd

import (
	"../cobbler"
	"../raft"
	"errors"
	"log"
)

// This command writes a value to a key.
type CobblerCommand struct {
	SerialNumber string                   `json:"serialNumber"`
	Action       string                   `json:"action"`
	Args         []map[string]interface{} `json:"args"`
}

// Creates a new write command.
func NewCobblerCommand(serialNumber string, action string, args []map[string]interface{}) *CobblerCommand {
	return &CobblerCommand{
		SerialNumber: serialNumber,
		Action:       action,
		Args:         args,
	}
}

// The name of the command in the log.
func (c *CobblerCommand) CommandName() string {
	return "cobbler"
}

func (c *CobblerCommand) Apply(server raft.Server) (interface{}, error) {
	cli := cobbler.NewClient()
	var err error
	switch c.Action {
	case "edit":
		for _, args := range c.Args {
			if err = cli.EditSystem(c.SerialNumber, args); err != nil {
				return nil, err
			}
		}
	case "add":
		if err = cli.AddSystem(c.SerialNumber, c.Args); err != nil {
			return nil, err
		}
	case "remove":
		if err = cli.RemoveSystem(c.SerialNumber); err != nil {
			return nil, err
		}
	default:
		log.Println("unsupported action")
		err = errors.New("unsupported action")
	}
	return nil, err
}
