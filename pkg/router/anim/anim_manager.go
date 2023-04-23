package anim

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
)

var (
	managers     = make(map[string]*anim.AnimManager)
	clientMgrMap = make(map[string]string) // Key: clientID, Value: managerID
)

func NewManager(clientIDs [2]string) string {
	id := uuid.New().String()
	managers[id] = anim.New()
	for i := 0; i < len(clientIDs); i++ {
		clientMgrMap[clientIDs[i]] = id
	}
	return id
}

func New(clientID string, a anim.Anim) string {
	mgrID := clientMgrMap[clientID]
	return managers[mgrID].New(a)
}

func Delete(clientID string, animID string) {
	mgrID := clientMgrMap[clientID]
	managers[mgrID].Delete(animID)
}

func GetAll(mgrID string) []anim.Param {
	return managers[mgrID].GetAll()
}

func MgrProcess(mgrID string) error {
	if mgrID == "" {
		return nil
	}

	return managers[mgrID].Process()
}
