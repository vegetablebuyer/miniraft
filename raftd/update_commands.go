package raftd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"../raft"
)

const (
	UpdateInitDir = "/var/log/cobbler/init_task/"
	UpdatePoolDir = "/var/log/cobbler/pool/"
)

// This command writes a value to a key.
type UpdateCommand struct {
	SerialNumber string `json:"serialNumber"`
	Name         string `json:"name"`
	Action       string `json:"action"`
	Info         string `json:"info"`
}

type UpdateResult struct {
	IsSucceed    bool   `json:"is_succeed"`
	SerialNumber string `json:"serial_number"`
	Result       string `json:"result"`
}

//system	819445715	192.168.10.73	1581415332.82
//29873	819188585	192.168.10.98	1567761609.17
// Creates a new write command.
func NewUpdateCommand(serialNumber string, name string, action string, typeD string, ipAddress string) *UpdateCommand {
	info := fmt.Sprintf("%v\t%v\t%v\t%v\n", typeD, serialNumber, ipAddress, time.Now().Unix())
	return &UpdateCommand{
		SerialNumber: serialNumber,
		Name:         name,
		Action:       action,
		Info:         info,
	}
}

// The name of the command in the log.
func (u *UpdateCommand) CommandName() string {
	return "update"
}

func (u *UpdateCommand) Apply(server raft.Server) (interface{}, error) {
	var err error
	action := u.Action
	switch action {
	case "pool":
		err = writeUpdateFile(filepath.Join(UpdatePoolDir, u.Name), os.O_CREATE|os.O_WRONLY, u.Info)
	case "task":
		err = writeUpdateFile(filepath.Join(UpdateInitDir, u.Name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, u.Info)
		if err == nil {
			_ = removePoolFile(u.SerialNumber)
		}
	default:
		log.Println("unsupported action")
		err = errors.New("unsupported action")
		return nil, err
	}

	return nil, err
}

func writeUpdateFile(fileName string, openFlag int, content string) error {
	f, err := os.OpenFile(fileName, openFlag, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

func removePoolFile(name string) (err error) {
	fileName := filepath.Join(UpdatePoolDir, name)
	if _, err = os.Stat(fileName); err == nil {
		err = os.Remove(fileName)
	}
	return
}
