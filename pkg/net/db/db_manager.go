package db

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db/memory"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db/model"
)

type Manager struct {
	client model.ClientHandler
	route  model.RouteHandler
}

var inst *Manager

func InitManager(dbType string, connStr string) error {
	if inst != nil {
		return fmt.Errorf("DB manager is already initialized")
	}

	switch dbType {
	case "memory":
		logger.Info("Initialize with local memory DB")
		inst = &Manager{
			client: memory.NewClientHandler(),
			route:  memory.NewRouteHandler(),
		}
	default:
		return fmt.Errorf("database type %s is not implemented yet", dbType)
	}

	return nil
}

func GetInst() *Manager {
	return inst
}

func (m *Manager) ClientAdd(ent model.ClientInfo) error {
	// TODO validation
	return m.client.Add(ent)
}

func (m *Manager) ClientDelete(clientID string) error {
	return m.client.Delete(clientID)
}

func (m *Manager) ClientGet() ([]model.ClientInfo, error) {
	return m.client.GetAll()
}

func (m *Manager) ClientGetByID(id string) (*model.ClientInfo, error) {
	// TODO refactoring
	clients, err := m.client.GetAll()
	if err != nil {
		return nil, err
	}
	for _, c := range clients {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("no such client")
}

func (m *Manager) RouteAdd(ent model.RouteInfo) error {
	// TODO validation
	return m.route.Add(ent)
}

func (m *Manager) RouteDelete(routeID string) error {
	route, err := m.route.Get(routeID)
	if err != nil {
		return fmt.Errorf("route get failed: %w", err)
	}

	for _, cid := range route.Clients {
		m.ClientDelete(cid)
	}

	return m.route.Delete(routeID)
}

func (m *Manager) RouteGet() ([]model.RouteInfo, error) {
	return m.route.GetAll()
}

func (m *Manager) RouteGetByClient(clientID string) (*model.RouteInfo, error) {
	// TODO refactoring
	routes, err := m.route.GetAll()
	if err != nil {
		return nil, err
	}
	for _, r := range routes {
		for _, c := range r.Clients {
			if c == clientID {
				return &r, nil
			}
		}
	}
	return nil, fmt.Errorf("no route for client %s", clientID)
}
