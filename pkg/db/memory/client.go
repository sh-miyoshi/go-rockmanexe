package memory

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/db/model"
)

type ClientInfoHandler struct {
	clientList []model.ClientInfo
}

func NewClientHandler() *ClientInfoHandler {
	res := &ClientInfoHandler{}
	return res
}

func (h *ClientInfoHandler) Add(ent model.ClientInfo) error {
	h.clientList = append(h.clientList, ent)
	return nil
}

func (h *ClientInfoHandler) Delete(clientID string) error {
	newList := []model.ClientInfo{}
	found := false
	for _, c := range h.clientList {
		if c.ID == clientID {
			found = true
		} else {
			newList = append(newList, c)
		}
	}

	if found {
		h.clientList = newList
		return nil
	}

	return fmt.Errorf("no such client %s", clientID)
}

func (h *ClientInfoHandler) Get() ([]model.ClientInfo, error) {
	return h.clientList, nil
}
