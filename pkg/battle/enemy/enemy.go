package enemy

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
)

// EnemyParam ...
type EnemyParam struct {
	CharID   int
	ObjectID string
	PlayerID string
	PosX     int
	PosY     int
	HP       int
}

type enemyObject interface {
	anim.Anim
	Init(ID string) error
	End()
}

var (
	ErrGameEnd = errors.New("game end")
	enemies    = make(map[string]enemyObject)
)

func Init(playerID string, enemyList []EnemyParam) error {
	for _, e := range enemyList {
		e.PlayerID = playerID
		obj := getObject(e.CharID, e)
		objID := anim.New(obj)
		enemies[objID] = obj
	}

	// Init enemy data
	for id, e := range enemies {
		if err := e.Init(id); err != nil {
			return err
		}
	}
	return nil
}

func End() {
	// Cleanup existsting enemy data
	for _, e := range enemies {
		e.End()
	}
}

func MgrProcess() error {
	for id, e := range enemies {
		if !anim.IsProcessing(id) {
			e.End()
			delete(enemies, id)
		}
	}

	if len(enemies) == 0 {
		return ErrGameEnd
	}

	return nil
}

func GetEnemyPositions() []field.ObjectPosition {
	res := []field.ObjectPosition{}
	for id, e := range enemies {
		pm := e.GetParam()
		res = append(res, field.ObjectPosition{
			ID: id,
			X:  pm.PosX,
			Y:  pm.PosY,
		})
	}
	return res
}
