package routerapi

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db/model"
)

func ClientGetV1(w http.ResponseWriter, r *http.Request) {
	clients, err := db.GetInst().ClientGet()
	if err != nil {
		logger.Error("Failed to get client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []ClientInfo{}
	for _, c := range clients {
		res = append(res, ClientInfo{
			ID:  c.ID,
			Key: c.Key,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("ClientGetV1 method successfully finished")
}

func ClientAddV1(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	key := uuid.New().String()
	ent := model.ClientInfo{
		ID:  id,
		Key: key,
	}

	if err := db.GetInst().ClientAdd(ent); err != nil {
		logger.Error("Failed to add client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := &ClientInfo{
		ID:  id,
		Key: key,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("ClientAddV1 method successfully finished")
}

func ClientDeleteV1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["clientID"]

	if err := db.GetInst().ClientDelete(clientID); err != nil {
		logger.Error("Failed to delete client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("ClientDeleteV1 method successfully finished")
}
