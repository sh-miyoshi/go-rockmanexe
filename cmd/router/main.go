package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconn"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gamehandler"
	"google.golang.org/grpc"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "config.yaml", "file path of config")
	flag.Parse()

	// Initialize config
	if err := config.Init(confFile); err != nil {
		fmt.Printf("Failed to init config: %v", err)
		os.Exit(1)
	}

	c := config.Get()
	logger.InitLogger(c.Log.DebugLog, c.Log.FileName)

	// Init Chip Info
	if err := chip.Init(c.ChipFilePath); err != nil {
		logger.Error("Failed to initialize chip info: %+v", err)
		return
	}

	fps.FPS = 60

	// Start API Server
	if c.Server.Enabled {
		logger.Info("start api server with %s", config.APIAddr())
		r := mux.NewRouter()
		setAPI(r)

		if err := http.ListenAndServe(config.APIAddr(), r); err != nil {
			logger.Error("Failed to run API server: %v", err)
			return
		}
	}

	// Listen data connection
	logger.Info("start data stream with %s", c.DataStreamAddr)
	listen, err := net.Listen("tcp", c.DataStreamAddr)
	if err != nil {
		logger.Error("Failed to listen data stream: %v", err)
		return
	}
	session.SetLogicGenerator(gamehandler.NewHandler)

	go session.ManagerExec()

	s := grpc.NewServer()
	pb.RegisterNetConnServer(s, netconn.New())

	if err = s.Serve(listen); err != nil {
		logger.Error("Failed to start data stream: %v", err)
		return
	}
}
