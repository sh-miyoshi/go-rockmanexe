package skillmanager

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Get(id int, arg skillcore.Argument) skillcore.SkillCore {
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
	case resources.SkillFlamePillarLine, resources.SkillFlamePillarTracking:
		res := &processor.FlamePillarManager{Arg: arg}
		res.Init(id)
		return res
	case resources.SkillTornado:
		res := &processor.Tornado{Arg: arg}
		res.Init()
		return res
	case resources.SkillBoomerang:
		res := &processor.Boomerang{Arg: arg}
		res.Init()
		return res
	case resources.SkillBambooLance:
		return &processor.BambooLance{Arg: arg}
	case resources.SkillCrackout, resources.SkillDoubleCrack, resources.SkillTripleCrack:
		res := &processor.Crack{Arg: arg}
		res.Init(id)
		return res
	default:
		// TODO: 不正なIDの場合はエラーをセットしたいが、現状実装途中なので呼び出し元で参照しないようにする
		logger.Error("skill %d is not implemented yet", id)
	}

	return nil
}
