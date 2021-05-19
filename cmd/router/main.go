package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/dstream"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

func main() {
	var exitErr chan error

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

	if err := db.InitManager(c.DB.Type, c.DB.ConnString); err != nil {
		logger.Error("Failed ot init DB manager: %v", err)
		os.Exit(1)
	}

	// Listen API request
	logger.Info("start API server with %s", c.APIAddr)
	go func() {
		r := mux.NewRouter()
		setAPI(r)

		if err := http.ListenAndServe(c.APIAddr, r); err != nil {
			logger.Error("Failed to run API server: %v", err)
			exitErr <- err
		}
	}()

	// Listen data connection
	logger.Info("start data stream with %s", c.DataStreamAddr)
	go func() {
		listen, err := net.Listen("tcp", c.DataStreamAddr)
		if err != nil {
			logger.Error("Failed to listen data stream: %v", err)
			exitErr <- err
		}

		s := grpc.NewServer()
		pb.RegisterRouterServer(s, &dstream.RouterStream{})

		if err = s.Serve(listen); err != nil {
			logger.Error("Failed to start data stream: %v", err)
			exitErr <- err
		}
	}()

	<-exitErr
}

func setAPI(r *mux.Router) {
	basePath := "/api/v1"

	r.HandleFunc(basePath+"/client", routerapi.ClientGetV1).Methods("GET")
	r.HandleFunc(basePath+"/client", routerapi.ClientAddV1).Methods("POST")
	r.HandleFunc(basePath+"/client/{clientID}", routerapi.ClientDeleteV1).Methods("DELETE")

	r.HandleFunc(basePath+"/route", routerapi.RouteGetV1).Methods("GET")
	r.HandleFunc(basePath+"/route", routerapi.RouteAddV1).Methods("POST")
	r.HandleFunc(basePath+"/route/{routeID}", routerapi.RouteDeleteV1).Methods("DELETE")
}
