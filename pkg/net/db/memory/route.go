package memory

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db/model"
)

type RouteInfoHandler struct {
	routeList []model.RouteInfo
}

func NewRouteHandler() *RouteInfoHandler {
	res := &RouteInfoHandler{}
	return res
}

func (h *RouteInfoHandler) Add(ent model.RouteInfo) error {
	h.routeList = append(h.routeList, ent)
	return nil
}

func (h *RouteInfoHandler) Delete(routeID string) error {
	newList := []model.RouteInfo{}
	found := false
	for _, c := range h.routeList {
		if c.ID == routeID {
			found = true
		} else {
			newList = append(newList, c)
		}
	}

	if found {
		h.routeList = newList
		return nil
	}

	return fmt.Errorf("no such route %s", routeID)
}

func (h *RouteInfoHandler) Get() ([]model.RouteInfo, error) {
	return h.routeList, nil
}
