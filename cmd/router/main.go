package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db/model"
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

	// Add debug clients and route
	if c.Debug.Enabled {
		db.GetInst().ClientAdd(model.ClientInfo{ID: c.Debug.Client1ID, Key: c.Debug.Client1Key})
		db.GetInst().ClientAdd(model.ClientInfo{ID: c.Debug.Client2ID, Key: c.Debug.Client2Key})
		route := model.RouteInfo{
			ID:      uuid.New().String(),
			Clients: [2]string{c.Debug.Client1ID, c.Debug.Client2ID},
		}
		db.GetInst().RouteAdd(route)
		logger.Info("Add debug clients")
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
		pb.RegisterRouterServer(s, dstream.New())

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

	r.Use(authMiddleware)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request
		authUser := os.Getenv("ROCKMAN_API_AUTH_USER")
		authPassword := os.Getenv("ROCKMAN_API_AUTH_PASSWORD")
		if len(authUser) > 0 && len(authPassword) > 0 {
			reqUser, reqPw, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if authUser != reqUser || authPassword != reqPw {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
