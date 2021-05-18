package db

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/db/memory"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/db/model"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
