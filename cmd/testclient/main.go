package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/app"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
)

const (
	streamAddr = "localhost:16283"
)

func main() {
	var clientID string
	var logfile string
	flag.StringVar(&clientID, "client", "", "client id")
	flag.StringVar(&clientID, "c", "", "client id")
	flag.StringVar(&logfile, "log", "", "log file")
	flag.Parse()

	if clientID == "" {
		fmt.Println("Please set client ID")
		return
	}

	if logfile != "" {
		file, err := os.Create(logfile)
		if err != nil {
			fmt.Printf("Failed to init logger: %v", err)
			return
		}
		log.SetOutput(file)
	}

	if err := app.Init(clientID); err != nil {
		log.Fatalf("Failed to init player info: %v", err)
		return
	}

	// run with debug client
	clientKey := "testtest"

	if err := netconn.Connect(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       clientID,
		ClientKey:      clientKey,
		ProgramVersion: "testclient",
	}); err != nil {
		log.Fatalf("Failed to connect router: %v", err)
		return
	}
	log.Println("Success to connect router")

	exitErr := make(chan error)
	go app.Process(exitErr)

	err := <-exitErr
	log.Fatalf("Run failed: %v", err)
}
