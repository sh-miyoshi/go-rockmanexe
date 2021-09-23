package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIAddr string `yaml:"api_addr"`
	Log     struct {
		FileName string `yaml:"file"`
	} `yaml:"log"`
	ClientData []struct {
		ClientID  string `yaml:"client_id"`
		ClientKey string `yaml:"client_key"`
		UserID    string `yaml:"user_id"`
		RouteID   string `json:"route_id"`
	} `yaml:"client_data"`
	SessionData []struct {
		ID            string `yaml:"id"`
		OwnerClientID string `yaml:"owner_client_id"`
		GuestClientID string `yaml:"guest_client_id"`
	} `yaml:"session_data"`
}

type AuthRequest struct {
	ClientID  string `json:"client_id"`
	ClientKey string `json:"client_key"`
}

type AuthResponse struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}

type SessionResponse struct {
	ID            string `json:"id"`
	OwnerClientID string `json:"owner_client_id"`
	GuestClientID string `json:"guest_client_id"`
}

var (
	cfg Config
)

func loadConfig(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer fp.Close()
	if err := yaml.NewDecoder(fp).Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode yaml: %v", err)
	}

	return nil
}

func setAPI(r *mux.Router) {
	r.HandleFunc("/api/v1/client/auth", func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode a request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		for _, d := range cfg.ClientData {
			if d.ClientID == req.ClientID && d.ClientKey == req.ClientKey {
				res := &AuthResponse{
					UserID:    d.UserID,
					SessionID: d.RouteID,
				}
				w.Header().Add("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(res); err != nil {
					logger.Error("Failed to encode a response: %+v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				logger.Info("Successfully auth user")

				return
			}
		}

		logger.Info("Failed to auth user with request: %+v", req)
		http.Error(w, "User Authentication Failed", http.StatusBadRequest)
	}).Methods("POST")

	r.HandleFunc("/api/v1/session/{sessionID}", func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}).Methods("GET")
}

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "config.yaml", "file path of config")
	flag.Parse()

	if err := loadConfig(confFile); err != nil {
		fmt.Printf("Failed to load config: %v", err)
		return
	}

	logger.InitLogger(true, cfg.Log.FileName)

	// Start API server
	r := mux.NewRouter()
	setAPI(r)

	logger.Info("start API server with %s", cfg.APIAddr)
	if err := http.ListenAndServe(cfg.APIAddr, r); err != nil {
		logger.Error("Failed to run API server: %v", err)
	}
}
