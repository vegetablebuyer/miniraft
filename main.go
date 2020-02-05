package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"./raft"
	"./raftd"
)

var host string
var port int
var join string

func init() {
	flag.StringVar(&host, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 4001, "port")
	flag.StringVar(&join, "join", "", "host:port a leader to join")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Data path argument is required")
	}

	rand.Seed(time.Now().Unix())

	raft.RegisterCommand(&raftd.WriteCommand{})
	raft.RegisterCommand(&raftd.CobblerCommand{})

	path := flag.Arg(0)
	if err := os.MkdirAll(path, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	s := raftd.NewServer(path, host, port)
	log.Fatal(s.ListenAndServe(join))
}
