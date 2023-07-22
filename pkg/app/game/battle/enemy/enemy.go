package enemy

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/stretchr/stew/slice"
)

const (
	IDMetall int = iota
	IDTarget
	IDBilly
	IDLark
	IDBoomer
	IDAquaman
	IDGaroo
	IDVolgear
	IDRockman
	IDSupportNPC
	IDColdman
)

type EnemyChipInfo struct {
	CharID        int
	ChipID        int
	Code          string
	RequiredLevel int
}

type EnemyParam struct {
	CharID          int
	ObjectID        string
	PlayerID        string
	Pos             common.Point
	HP              int
	ActNo           int
	InvincibleCount int
	DamageType      int
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
		{CharID: IDBoomer, ChipID: chip.IDBoomerang1, Code: "m", RequiredLevel: 7},
		{CharID: IDBoomer, ChipID: chip.IDBoomerang1, Code: "*", RequiredLevel: 9},
		{CharID: IDAquaman, ChipID: chip.IDAquaman, Code: "a", RequiredLevel: 9},
		{CharID: IDVolgear, ChipID: chip.IDFlameLine1, Code: "f", RequiredLevel: 7},
		{CharID: IDGaroo, ChipID: chip.IDHeatShot, Code: "c", RequiredLevel: 7},
		// TODO: コールドマンのチップ
	}
)

func Init(playerID string, enemyList []EnemyParam) error {
	for i, e := range enemyList {
		e.PlayerID = playerID
		e.ActNo = i
		obj := getObject(e.CharID, e)
		objID := localanim.ObjAnimNew(obj)
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
		if !localanim.ObjAnimIsProcessing(id) {
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
	name = path + GetName(id)
	return
}

func GetName(id int) string {
	switch id {
	case IDMetall:
		return "メットール"
	case IDTarget:
		return "的"
	case IDBilly:
		return "ビリー"
	case IDLark:
		return "ゲイラーク"
	case IDBoomer:
		return "ラウンダ"
	case IDAquaman:
		return "アクアマン"
	case IDGaroo:
		return "ガルー"
	case IDVolgear:
		return "ボルケギア"
	case IDRockman:
		return "ロックマン"
	case IDColdman:
		return "コールドマン"
	}
	return ""
}

func IsBoss(id int) bool {
	bossList := []int{IDAquaman, IDRockman}
	return slice.Contains(bossList, id)
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
	case IDGaroo:
		return &enemyGaroo{pm: initParam}
	case IDVolgear:
		return &enemyVolgear{pm: initParam}
	case IDRockman:
		common.SetError("enemy rockman is not implemented yet")
	case IDColdman:
		return &enemyColdman{pm: initParam}
	}
	return nil
}

func damageProc(dm *damage.Damage, pm *EnemyParam) bool {
	if dm == nil {
		return false
	}

	if dm.TargetType&damage.TargetEnemy != 0 {
		if pm.InvincibleCount > 0 && dm.Power > 0 {
			return false
		}

		if damage.IsWeakness(0, *dm) {
			dm.Power *= 2
			localanim.AnimNew(effect.Get(resources.EffectTypeExclamation, pm.Pos, 0))
		}

		pm.HP -= dm.Power

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&pm.Pos, common.DirectLeft, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&pm.Pos, common.DirectRight, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}

		localanim.AnimNew(effect.Get(dm.HitEffectType, pm.Pos, 5))
		return true
	}
	return false
}

/*
Enemy template

package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
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
	return damageProc(dm, &e.pm)
}

func (e *enemy) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemy) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemy) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

*/
