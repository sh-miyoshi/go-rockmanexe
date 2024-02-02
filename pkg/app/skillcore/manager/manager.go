package skillmanager

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Manager struct {
	damageMgr    *damage.DamageManager
	GetObjectPos func(objID string) point.Point
}

func NewManager(damageMgr *damage.DamageManager, GetObjectPos func(objID string) point.Point) *Manager {
	return &Manager{
		damageMgr:    damageMgr,
		GetObjectPos: GetObjectPos,
	}
}

func (m *Manager) Get(id int, arg skillcore.Argument) skillcore.SkillCore {
	switch id {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return &processor.Cannon{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos}
	case resources.SkillMiniBomb:
		return &processor.MiniBomb{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos}
	case resources.SkillRecover:
		return &processor.Recover{Arg: arg, DamageMgr: m.damageMgr}
	case resources.SkillEnemyShockWave:
		res := &processor.ShockWave{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos}
		res.Init(false)
		return res
	case resources.SkillPlayerShockWave:
		res := &processor.ShockWave{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos}
		res.Init(true)
		return res
	case resources.SkillSpreadGun:
		return &processor.SpreadGun{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos}
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword:
		return &processor.Sword{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos, SkillID: id}
	case resources.SkillVulcan1:
		return &processor.Vulcan{Arg: arg, DamageMgr: m.damageMgr, GetObjectPos: m.GetObjectPos, Times: 3}
	}

	// TODO: 不正なIDの場合はエラーをセットしたいが、現状実装途中なので呼び出し元で参照しないようにする
	return nil
}
