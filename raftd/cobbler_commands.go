package raftd

import (
	"../cobbler"
	"../raft"
	"errors"
	"log"
)

// This command writes a value to a key.
type CobblerCommand struct {
	SerialNumber string      `json:"serialNumber"`
	Action       string      `json:"action"`
	Args         interface{} `json:"args"`
}

type CobblerResult struct {
	IsSucceed    bool   `json:"is_succeed"`
	SerialNumber string `json:"serial_number"`
	Result       string `json:"result"`
}

// Creates a new write command.
func NewCobblerCommand(serialNumber string, action string, args interface{}) *CobblerCommand {
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
		value, ok := c.Args.([]interface{})
		if !ok {
			err = errors.New("system edit args should be []map[string]interface{}")
		}
		for _, args := range value {
			if err = cli.EditSystem(c.SerialNumber, args); err != nil {
				return nil, err
			}
		}
	case "add":
		value, ok := c.Args.(map[string]interface{})
		if !ok {
			err = errors.New("system add args should be map[string]interface{}")
		}
		if err = cli.AddSystem(c.SerialNumber, value); err != nil {
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
