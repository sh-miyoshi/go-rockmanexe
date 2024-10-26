package skillmanager

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
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
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword, resources.SkillFighterSword, resources.SkillNonEffectWideSword:
		return &processor.Sword{Arg: arg, SkillID: id}
	case resources.SkillVulcan1:
		return &processor.Vulcan{Arg: arg, Times: 3}
	case resources.SkillVulcan2:
		return &processor.Vulcan{Arg: arg, Times: 5}
	case resources.SkillVulcan3:
		return &processor.Vulcan{Arg: arg, Times: 7}
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
	case resources.SkillCountBomb:
		return &processor.CountBomb{Arg: arg}
	case resources.SkillAreaSteal, resources.SkillPanelSteal:
		res := &processor.AreaSteal{Arg: arg}
		res.Init(id)
		return res
	case resources.SkillAquaman:
		res := &processor.Aquaman{Arg: arg}
		res.Init()
		return res
	case resources.SkillInvisible:
		return &processor.Invisible{Arg: arg}
	case resources.SkillQuickGauge:
		return &processor.QuickGauge{Arg: arg}
	case resources.SkillAquamanShot, resources.SkillCirkillShot, resources.SkillGarooBreath, resources.SkillShrimpyAttack:
		// 敵の攻撃の場合、共通化するメリットがないのでcoreは何もしない
		return nil
	case resources.SkillThunderBall:
		res := &processor.ThunderBall{Arg: arg}
		res.Init()
		return res
	case resources.SkillBubbleShot, resources.SkillBubbleV, resources.SkillBubbleSide:
		return &processor.BubbleShot{Arg: arg, SkillID: id}
	case resources.SkillForteHellsRollingUp, resources.SkillForteHellsRollingDown:
		res := &processor.ForteHellsRolling{Arg: arg}
		res.Init(id, false)
		return res
	case resources.SkillForteDarkArmBladeType1, resources.SkillForteDarkArmBladeType2:
		res := &processor.ForteDarkArmBlade{Arg: arg}
		res.Init(id)
		return res
	case resources.SkillForteShootingBuster:
		res := &processor.ForteShootingBuster{Arg: arg}
		res.Init()
		return res
	case resources.SkillForteDarknessOverload:
		return &processor.ForteDarknessOverload{Arg: arg}
	case resources.SkillChipForteAnother:
		res := &processor.ChipForteAnother{Arg: arg}
		res.Init()
		return res
	case resources.SkillDeathMatch1, resources.SkillDeathMatch2:
		res := &processor.DeathMatch{Arg: arg, SkillID: id}
		res.Init()
		return res
	case resources.SkillFullCustom:
		return &processor.FullCustom{Arg: arg}
	case resources.SkillPanelReturn:
		return &processor.PanelReturn{Arg: arg}
	default:
		system.SetError(fmt.Sprintf("skill %d is not implemented yet", id))
	}

	return nil
}
