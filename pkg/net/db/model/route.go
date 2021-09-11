package model

import "errors"

var (
	ErrNoSuchRoute = errors.New("no such route")
)

type RouteInfo struct {
	ID      string
	Clients [2]string
}

type RouteHandler interface {
	Add(ent RouteInfo) error
	Delete(routeID string) error
	GetAll() ([]RouteInfo, error)
	Get(routeID string) (*RouteInfo, error)
}
