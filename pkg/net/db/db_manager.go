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
	return m.client.Get()
}

func (m *Manager) ClientGetByID(id string) (*model.ClientInfo, error) {
	// TODO refactoring
	clients, err := m.client.Get()
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
	return m.route.Delete(routeID)
}

func (m *Manager) RouteGet() ([]model.RouteInfo, error) {
	return m.route.Get()
}
