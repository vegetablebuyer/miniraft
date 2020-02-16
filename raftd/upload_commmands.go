package raftd

import (
	"../cobbler"
	"../raft"
	"log"
	"os"
	"path/filepath"
)

const UploadDir = "/var/log/cobbler/report/"

// This command writes a value to a key.
type UploadCommand struct {
	SerialNumber string `json:"serialNumber"`
	Info         string `json:"info"`
}

type UploadResult struct {
	IsSucceed    bool   `json:"is_succeed"`
	SerialNumber string `json:"serial_number"`
	Result       string `json:"result"`
}

// Creates a new write command.
func NewUploadCommand(serialNumber string, info string) *UploadCommand {
	return &UploadCommand{
		SerialNumber: serialNumber,
		Info:         info,
	}
}

// The name of the command in the log.
func (u *UploadCommand) CommandName() string {
	return "upload"
}

func (u *UploadCommand) Apply(server raft.Server) (interface{}, error) {
	log.Println("start to upload")
	err := writeReportFile(u.SerialNumber, u.Info)
	cli := cobbler.NewClient()
	if !cli.FindSystem(u.SerialNumber) {
		return nil, nil
	}
	arg := make(map[string]interface{})
	ksMeta := make(map[string]string)
	ksMeta["report"] = "0"
	arg["ksmeta"] = ksMeta
	if err = cli.EditSystem(u.SerialNumber, arg); err != nil {
		return nil, err
	}
	return nil, err
}

func writeReportFile(name, content string) error {
	fileName := filepath.Join(UploadDir, name)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}
