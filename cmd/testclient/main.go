package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/app"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
)

const (
	streamAddr = "localhost:80"
)

func main() {
	var clientID string
	flag.StringVar(&clientID, "client", "", "client id")
	flag.StringVar(&clientID, "c", "", "client id")
	flag.Parse()

	if clientID == "" {
		fmt.Println("Please set client ID")
		return
	}

	if err := app.PlayerInit(); err != nil {
		log.Fatalf("Failed to init player info: %v", err)
		return
	}

	// run with debug client
	clientKey := "testtest"

	if err := netconn.Connect(streamAddr, clientID, clientKey); err != nil {
		log.Fatalf("Failed to connect router: %v", err)
		return
	}
	log.Println("Success to connect router")

	exitErr := make(chan error)
	go app.PlayerProc(exitErr)

	err := <-exitErr
	log.Fatalf("Run failed: %v", err)
}
