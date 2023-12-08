package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

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
		SessionID string `yaml:"session_id"`
	} `yaml:"client_data"`
	SessionData []struct {
		ID            string `yaml:"id"`
		OwnerUserID   string `yaml:"owner_user_id"`
		OwnerClientID string `yaml:"owner_client_id"`
		GuestUserID   string `yaml:"guest_user_id"`
		GuestClientID string `yaml:"guest_client_id"`
	} `yaml:"session_data"`
}

type AuthRequest struct {
	ClientID  string `json:"client_id"`
	ClientKey string `json:"client_key"`
}

type AuthResponse struct {
	SessionID string `json:"session_id"`
}

type SessionResponse struct {
	ID            string `json:"id"`
	OwnerUserID   string `json:"owner_user_id"`
	OwnerClientID string `json:"owner_client_id"`
	GuestUserID   string `json:"guest_user_id"`
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
					SessionID: d.SessionID,
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
		vars := mux.Vars(r)
		sessionID := vars["sessionID"]

		for _, d := range cfg.SessionData {
			if d.ID == sessionID {
				res := &SessionResponse{
					ID:            d.ID,
					OwnerUserID:   d.OwnerUserID,
					OwnerClientID: d.OwnerClientID,
					GuestUserID:   d.GuestUserID,
					GuestClientID: d.GuestClientID,
				}
				w.Header().Add("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(res); err != nil {
					logger.Error("Failed to encode a response: %+v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				logger.Info("Successfully return session info")
				return
			}
		}

		logger.Info("Failed to get session for: %s", sessionID)
		http.Error(w, "No such session", http.StatusNotFound)
	}).Methods("GET")

	r.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		reqBuf := new(bytes.Buffer)
		reqBuf.ReadFrom(r.Body)
		logger.Info("Request data: %+v", reqBuf.String())

		sec := 5 * time.Second
		logger.Info("Waiting %d[sec] ...", sec)
		time.Sleep(sec)

		w.Header().Add("Content-Type", "application/json")

		res := `{
	"id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-3.5-turbo-0613",
  "system_fingerprint": "fp_44709d6fcb",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "テストメッセージ"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}`

		w.Write([]byte(res))
	}).Methods("POST")
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
	logger.Debug("config: %+v", cfg)

	// Start API server
	r := mux.NewRouter()
	setAPI(r)

	logger.Info("start API server with %s", cfg.APIAddr)
	if err := http.ListenAndServe(cfg.APIAddr, r); err != nil {
		logger.Error("Failed to run API server: %v", err)
	}
}
