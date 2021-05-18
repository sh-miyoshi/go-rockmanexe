package model

type RouteInfo struct {
	ID      string
	Clients [2]string
}

type RouteHandler interface {
	Add(ent RouteInfo) error
	Delete(routeID string) error
	Get() ([]RouteInfo, error)
}
