package enemy

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
)

type EnemyChipInfo struct {
	CharID        int
	ChipID        int
	Code          string
	RequiredLevel int
}

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

	enemyChipList = []EnemyChipInfo{
		{CharID: IDMetall, ChipID: chip.IDShockWave, Code: "l", RequiredLevel: 7},
	}
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

func GetEnemyChip(id int, bustingLv int) []EnemyChipInfo {
	res := []EnemyChipInfo{}
	for _, c := range enemyChipList {
		if c.CharID == id && bustingLv >= c.RequiredLevel {
			res = append(res, c)
		}
	}
	return res
}
