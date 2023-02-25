package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/botclient/app"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
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

	clientKey := "testtest"
	connInst := netconn.New(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       clientID,
		ClientKey:      clientKey,
		ProgramVersion: "botclient",
		Insecure:       true,
	})
	connInst.ConnectRequest()

	app.Init(clientID, connInst)

	// Waiting connection
	for {
		st := connInst.GetConnStatus()
		if st.Status == netconn.ConnStateOK {
			break
		}
		if st.Status == netconn.ConnStateError {
			logger.Error("Failed to connect router: %v", st.Error)
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("Success to connect router")

	if err := app.Process(); err != nil {
		logger.Error("Failed to run app: %v", err)
	}

	// 終了処理をできるように若干待つ
	time.Sleep(1 * time.Second)
	logger.Info("Successfully closed app")
}
