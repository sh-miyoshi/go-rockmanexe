package main

import (
	"flag"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/router/config"
	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/api/router"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "config.yaml", "file path of config")
	flag.Parse()

	// initialize config
	if err := config.Init(confFile); err != nil {
		log.Fatalf("Failed to init config: %v", err)
		os.Exit(1)
	}

	// set API
	r := mux.NewRouter()
	setAPI(r)

	// listen data connection
}

func setAPI(r *mux.Router) {
	basePath := "api/v1"

	r.HandleFunc(basePath+"/client", routerapi.ClientAddV1).Methods("POST")
}
