package enemy

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
)

const (
	IDMetall int = iota
	IDTarget
	IDBilly
	IDLark
	IDBoomer
	IDAquaman
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
	ActNo    int
}

type enemyObject interface {
	objanim.Anim
	Init(ID string) error
	End()
}

var (
	ErrGameEnd = errors.New("game end")
	enemies    = make(map[string]enemyObject)

	enemyChipList = []EnemyChipInfo{
		{CharID: IDMetall, ChipID: chip.IDShockWave, Code: "l", RequiredLevel: 7},
		{CharID: IDMetall, ChipID: chip.IDShockWave, Code: "*", RequiredLevel: 9},
		{CharID: IDBilly, ChipID: chip.IDThunderBall, Code: "l", RequiredLevel: 7},
		{CharID: IDLark, ChipID: chip.IDWideShot, Code: "c", RequiredLevel: 7},
		// TODO boomer chip
		// TODO aquaman chip
	}
)

func Init(playerID string, enemyList []EnemyParam) error {
	for i, e := range enemyList {
		e.PlayerID = playerID
		e.ActNo = i
		obj := getObject(e.CharID, e)
		objID := objanim.New(obj)
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
		if !objanim.IsProcessing(id) {
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
	case IDBoomer:
		name = path + "ラウンダ"
	case IDAquaman:
		name = path + "アクアマン"
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
	case IDBoomer:
		return &enemyBoomer{pm: initParam}
	case IDAquaman:
		return &enemyAquaman{pm: initParam}
	}
	return nil
}

/*
Enemy template

package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
)

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

func (e *enemy) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
		return true
	}
	return false
}

func (e *enemy) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemy) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

*/
