package enemy

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

const (
	IDMetall int = iota
	IDTarget
	IDBilly
	IDLark
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

func GetStandImageFile(id int) (name, ext string) {
	ext = ".png"

	path := common.ImagePath + "battle/character/"

	switch id {
	case IDMetall:
		name = path + "メットール"
	case IDTarget:
		name = path + "的"
	case IDBilly:
		name = path + "ビリー"
	case IDLark:
		name = path + "ゲイラーク"
	}
	return
}

func getObject(id int, initParam EnemyParam) enemyObject {
	switch id {
	case IDMetall:
		return &enemyMetall{pm: initParam}
	case IDTarget:
		return &enemyTarget{pm: initParam}
	case IDBilly:
		return &enemyBilly{pm: initParam}
	case IDLark:
		return &enemyLark{pm: initParam}
	}
	return nil
}

/*
Enemy template

type enemy struct {
	pm EnemyParam
}

func (e *enemy) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	return nil
}

func (e *enemy) End() {
	// Delete Images
}

func (e *enemy) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemy) Draw() {
	// Show Enemy Images
}

func (e *enemy) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
	}
}

func (e *enemy) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

*/
