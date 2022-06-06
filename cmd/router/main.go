package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconn"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
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

	// Listen data connection
	logger.Info("start data stream with %s", c.DataStreamAddr)
	listen, err := net.Listen("tcp", c.DataStreamAddr)
	if err != nil {
		logger.Error("Failed to listen data stream: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterNetConnServer(s, netconn.New())

	if err = s.Serve(listen); err != nil {
		logger.Error("Failed to start data stream: %v", err)
		return
	}
}
