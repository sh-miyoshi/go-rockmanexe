package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/app"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

	logger.InitLogger(true, logfile)

	if err := app.Init(clientID); err != nil {
		logger.Error("Failed to init player info: %v", err)
		return
	}

	// run with debug client
	clientKey := "testtest"

	netconn.Connect(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       clientID,
		ClientKey:      clientKey,
		ProgramVersion: "testclient",
		Insecure:       true,
	})

	for {
		if err := netconn.GetConnectStatus(); err != nil {
			if err == netconn.ErrConnected {
				break
			} else {
				logger.Error("Failed to connect router: %v", err)
				return
			}
		}

		time.Sleep(100 * time.Microsecond)
	}

	logger.Info("Success to connect router")

	exitErr := make(chan error)
	go app.Process(exitErr)

	err := <-exitErr
	logger.Error("Run failed: %v", err)
}
