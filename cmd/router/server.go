package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/api"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
)

func setAPI(r *mux.Router) {
	r.HandleFunc("/api/v1/client/auth", func(w http.ResponseWriter, r *http.Request) {
		var req api.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode a request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		c := config.Get()
		ok := false
		if c.Server.Session.ClientID1 == req.ClientID {
			ok = c.Server.Session.ClientKey1 == req.ClientKey
		}
		if c.Server.Session.ClientID2 == req.ClientID {
			ok = c.Server.Session.ClientKey2 == req.ClientKey
		}

		if ok {
			res := &api.AuthResponse{
				SessionID: c.Server.Session.ID,
			}
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				logger.Error("Failed to encode a response: %+v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			logger.Info("Failed to auth user with request: %+v", req)
			http.Error(w, "User Authentication Failed", http.StatusBadRequest)
		}
	}).Methods("POST")

	r.HandleFunc("/api/v1/session/{sessionID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionID := vars["sessionID"]

		c := config.Get()
		if sessionID == c.Server.Session.ID {
			res := &api.SessionInfo{
				ID:            c.Server.Session.ID,
				OwnerUserID:   c.Server.Session.ClientID1,
				OwnerClientID: c.Server.Session.ClientID1,
				GuestUserID:   c.Server.Session.ClientID2,
				GuestClientID: c.Server.Session.ClientID2,
			}
			w.Header().Add("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				logger.Error("Failed to encode a response: %+v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			logger.Info("Successfully return session info")
		} else {
			logger.Info("Failed to get session for: %s", sessionID)
			http.Error(w, "No such session", http.StatusNotFound)
		}
	}).Methods("GET")
}
