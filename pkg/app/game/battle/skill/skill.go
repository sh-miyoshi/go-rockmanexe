package skill

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
)

type SkillAnim interface {
	anim.Anim

	StopByOwner()
}

func Init() error {
	if err := skilldraw.LoadImages(); err != nil {
		return errors.Wrap(err, "failed to load skill image")
	}

	return nil
}

func End() {
	skilldraw.ClearImages()
}

func Get(skillID int, arg skillcore.Argument) SkillAnim {
	objID := uuid.New().String()
	arg.GetPanelInfo = field.GetPanelInfo
	arg.ChangePanelStatus = field.ChangePanelStatus
	arg.DamageMgr = localanim.DamageManager()
	arg.GetObjectPos = localanim.ObjAnimGetObjPos
	arg.GetObjects = localanim.ObjAnimGetObjs
	arg.SoundOn = sound.On
	arg.Cutin = func(skillName string, count int) {
		field.SetBlackoutCount(count)
		SetChipNameDraw(skillName, true)
	}
	arg.MakeInvisible = localanim.ObjAnimMakeInvisible
	arg.AddBarrier = localanim.ObjAnimAddBarrier
	arg.ChangePanelType = field.ChangePanelType
	core := localanim.SkillManager().Get(skillID, arg)

	switch skillID {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(objID, arg, core, skillID)
	case resources.SkillMiniBomb:
		return newMiniBomb(objID, arg, core)
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword, resources.SkillFighterSword, resources.SkillNonEffectWideSword:
		return newSword(objID, arg, core)
	case resources.SkillPlayerShockWave, resources.SkillEnemyShockWave:
		return newShockWave(objID, arg, core)
	case resources.SkillRecover:
		return newRecover(objID, arg, core)
	case resources.SkillSpreadGun:
		return newSpreadGun(objID, arg, core)
	case resources.SkillVulcan1, resources.SkillVulcan2, resources.SkillVulcan3:
		return newVulcan(objID, arg, core)
	case resources.SkillThunderBall:
		return newThunderBall(objID, arg, core)
	case resources.SkillPlayerWideShot, resources.SkillEnemyWideShot:
		return newWideShot(objID, arg, core)
	case resources.SkillBoomerang:
		return newBoomerang(objID, arg, core)
	case resources.SkillWaterBomb:
		return newWaterBomb(objID, arg, core)
	case resources.SkillAquamanShot:
		return newAquamanShot(objID, arg)
	case resources.SkillAquaman:
		return newAquaman(objID, arg, core)
	case resources.SkillCrackout, resources.SkillDoubleCrack, resources.SkillTripleCrack:
		return newCrack(objID, arg, core)
	case resources.SkillBambooLance:
		return newBambooLance(objID, arg, core)
	case resources.SkillInvisible:
		return newInvisible(objID, arg, core)
	case resources.SkillGarooBreath:
		return newGarooBreath(objID, arg)
	case resources.SkillFlamePillarTracking, resources.SkillFlamePillarRandom, resources.SkillFlamePillarLine:
		return newFlamePillar(objID, arg, core)
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return newHeatShot(objID, arg, core)
	case resources.SkillAreaSteal, resources.SkillPanelSteal:
		return newAreaSteal(objID, arg, core)
	case resources.SkillCountBomb:
		return newCountBomb(objID, arg, core)
	case resources.SkillTornado:
		return newTornado(objID, arg, core)
	case resources.SkillFailed:
		return newFailed(objID, arg)
	case resources.SkillQuickGauge:
		return newQuickGauge(objID, arg, core)
	case resources.SkillCirkillShot:
		return newCirkillShot(objID, arg)
	case resources.SkillShrimpyAttack:
		return newShrimpyAtk(objID, arg)
	case resources.SkillBubbleShot, resources.SkillBubbleV, resources.SkillBubbleSide:
		return newBubbleShot(objID, arg, core)
	case resources.SkillForteHellsRollingUp, resources.SkillForteHellsRollingDown:
		return newForteHellsRolling(objID, arg, core)
	case resources.SkillForteDarkArmBladeType1, resources.SkillForteDarkArmBladeType2:
		return newForteDarkArmBlade(objID, arg, core, skillID)
	case resources.SkillForteShootingBuster:
		return newForteShootingBuster(objID, arg, core)
	case resources.SkillForteDarknessOverload:
		return newForteDarknessOverload(objID, arg, core)
	case resources.SkillChipForteAnother:
		return newChipForteAnother(objID, arg, core)
	case resources.SkillDeathMatch1, resources.SkillDeathMatch2, resources.SkillDeathMatch3:
		return newDeathMatch(objID, arg, core)
	case resources.SkillPanelReturn:
		return newPanelReturn(objID, arg, core)
	case resources.SkillBarrier, resources.SkillBarrier100, resources.SkillBarrier200:
		return newBarrier(objID, arg, core)
	}

	system.SetError(fmt.Sprintf("Skill %d is not implemented yet", skillID))
	return nil
}

/*
Skill template
package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type tmpskill struct {
	ID  string
	Arg skillcore.Argument
	Core *processor.TmpSkill
}

func newTmpSkill(objID string, arg skillcore.Argument, core skillcore.SkillCore) *tmpskill {
	return &tmpskill{
		ID:  objID,
		Arg: arg,
		Core: core.(*processor.TmpSkill),
	}
}

func (p *tmpskill) Draw() {
	// p.drawer.Draw()
}

func (p *tmpskill) Update() (bool, error) {
	return p.Core.Update()
}

func (p *tmpskill) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *tmpskill) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
*/
