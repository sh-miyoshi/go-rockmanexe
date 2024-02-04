package skillmanager

import (
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Manager struct {
	damageMgr    *damage.DamageManager
	GetObjectPos func(objID string) point.Point
	GetObjects   func(filter objanim.Filter) []objanim.Param
	SoundOn      func(typ resources.SEType)
}

func NewManager(
	damageMgr *damage.DamageManager,
	GetObjectPos func(objID string) point.Point,
	GetObjects func(filter objanim.Filter) []objanim.Param,
	SoundOn func(typ resources.SEType),
) *Manager {
	return &Manager{
		damageMgr:    damageMgr,
		GetObjectPos: GetObjectPos,
		GetObjects:   GetObjects,
		SoundOn:      SoundOn,
	}
}

func (m *Manager) Get(id int, arg skillcore.Argument) skillcore.SkillCore {
	arg.GetObjectPos = m.GetObjectPos
	arg.GetObjects = m.GetObjects
	arg.SoundOn = m.SoundOn
	arg.DamageMgr = m.damageMgr

	switch id {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return &processor.Cannon{SkillID: id, Arg: arg}
	case resources.SkillMiniBomb:
		res := &processor.MiniBomb{Arg: arg}
		res.Init()
		return res
	case resources.SkillRecover:
		return &processor.Recover{Arg: arg}
	case resources.SkillPlayerShockWave, resources.SkillEnemyShockWave:
		res := &processor.ShockWave{Arg: arg}
		res.Init(id == resources.SkillPlayerShockWave)
		return res
	case resources.SkillSpreadGun:
		return &processor.SpreadGun{Arg: arg}
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword:
		return &processor.Sword{Arg: arg, SkillID: id}
	case resources.SkillVulcan1:
		return &processor.Vulcan{Arg: arg, Times: 3}
	case resources.SkillPlayerWideShot, resources.SkillEnemyWideShot:
		res := &processor.WideShot{Arg: arg}
		res.Init(id == resources.SkillPlayerWideShot)
		return res
	case resources.SkillWaterBomb:
		res := &processor.WaterBomb{Arg: arg}
		res.Init()
		return res
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return &processor.HeatShot{Arg: arg, SkillID: id}
	}

	// TODO: 不正なIDの場合はエラーをセットしたいが、現状実装途中なので呼び出し元で参照しないようにする
	return nil
}
