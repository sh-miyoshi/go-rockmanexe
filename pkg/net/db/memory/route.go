package memory

import (
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
	for _, r := range h.routeList {
		if r.ID == routeID {
			found = true
		} else {
			newList = append(newList, r)
		}
	}

	if found {
		h.routeList = newList
		return nil
	}

	return model.ErrNoSuchRoute
}

func (h *RouteInfoHandler) GetAll() ([]model.RouteInfo, error) {
	return h.routeList, nil
}

func (h *RouteInfoHandler) Get(routeID string) (*model.RouteInfo, error) {
	for _, r := range h.routeList {
		if r.ID == routeID {
			return &r, nil
		}
	}

	return nil, model.ErrNoSuchRoute
}
