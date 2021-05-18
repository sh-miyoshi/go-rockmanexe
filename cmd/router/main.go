package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "config.yaml", "file path of config")
	flag.Parse()

	// initialize config
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

	// set API
	r := mux.NewRouter()
	setAPI(r)

	// listen data connection
}

func setAPI(r *mux.Router) {
	basePath := "api/v1"

	r.HandleFunc(basePath+"/client", routerapi.ClientAddV1).Methods("POST")
}
