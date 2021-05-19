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

func RouteGetV1(w http.ResponseWriter, r *http.Request) {
	routes, err := db.GetInst().RouteGet()
	if err != nil {
		logger.Error("Failed to get route: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []RouteInfo{}
	for _, r := range routes {
		res = append(res, RouteInfo{
			ID:      r.ID,
			Clients: r.Clients,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("RouteGetV1 method successfully finished")
}

func RouteAddV1(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()

	var req RouteAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode a request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ent := model.RouteInfo{
		ID:      id,
		Clients: req.Clients,
	}

	if err := db.GetInst().RouteAdd(ent); err != nil {
		logger.Error("Failed to add route: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := &RouteInfo{
		ID:      id,
		Clients: req.Clients,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("RouteAddV1 method successfully finished")
}

func RouteDeleteV1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeID := vars["routeID"]

	if err := db.GetInst().RouteDelete(routeID); err != nil {
		logger.Error("Failed to delete route: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("RouteDeleteV1 method successfully finished")
}
