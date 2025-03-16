package skill

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skillanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
)

type SkillAnim interface {
	skillanim.Anim

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

func Get(skillID int, arg skillcore.Argument, animMgr *manager.Manager) SkillAnim {
	objID := uuid.New().String()
	arg.GetPanelInfo = field.GetPanelInfo
	arg.ChangePanelStatus = field.ChangePanelStatus
	arg.DamageMgr = animMgr.DamageManager()
	arg.GetObjectPos = animMgr.ObjAnimGetObjPos
	arg.GetObjects = animMgr.ObjAnimGetObjs
	arg.SoundOn = sound.On
	arg.Cutin = func(skillName string, count int) {
		field.SetBlackoutCount(count)
		animMgr.SetActiveAnim(objID)
		SetChipNameDraw(skillName, true)
	}
	arg.MakeInvisible = animMgr.ObjAnimMakeInvisible
	arg.AddBarrier = animMgr.ObjAnimAddBarrier
	arg.ChangePanelType = field.ChangePanelType
	core := animMgr.SkillGet(skillID, arg)

	switch skillID {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(objID, arg, core, skillID, animMgr)
	case resources.SkillMiniBomb:
		return newMiniBomb(objID, arg, core, animMgr)
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword, resources.SkillFighterSword, resources.SkillNonEffectWideSword:
		return newSword(objID, arg, core, animMgr)
	case resources.SkillPlayerShockWave, resources.SkillEnemyShockWave:
		return newShockWave(objID, arg, core, animMgr)
	case resources.SkillRecover:
		return newRecover(objID, arg, core, animMgr)
	case resources.SkillSpreadGun:
		return newSpreadGun(objID, arg, core, animMgr)
	case resources.SkillVulcan1, resources.SkillVulcan2, resources.SkillVulcan3:
		return newVulcan(objID, arg, core, animMgr)
	case resources.SkillThunderBall:
		return newThunderBall(objID, arg, core, animMgr)
	case resources.SkillPlayerWideShot, resources.SkillEnemyWideShot:
		return newWideShot(objID, arg, core, animMgr)
	case resources.SkillBoomerang:
		return newBoomerang(objID, arg, core, animMgr)
	case resources.SkillWaterBomb:
		return newWaterBomb(objID, arg, core, animMgr)
	case resources.SkillAquamanShot:
		return newAquamanShot(objID, arg, animMgr)
	case resources.SkillAquaman:
		return newAquaman(objID, arg, core, animMgr)
	case resources.SkillCrackout, resources.SkillDoubleCrack, resources.SkillTripleCrack:
		return newCrack(objID, arg, core, animMgr)
	case resources.SkillBambooLance:
		return newBambooLance(objID, arg, core, animMgr)
	case resources.SkillInvisible:
		return newInvisible(objID, arg, core, animMgr)
	case resources.SkillGarooBreath:
		return newGarooBreath(objID, arg, animMgr)
	case resources.SkillFlamePillarTracking, resources.SkillFlamePillarRandom, resources.SkillFlamePillarLine:
		return newFlamePillar(objID, arg, core, animMgr)
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return newHeatShot(objID, arg, core, animMgr)
	case resources.SkillAreaSteal, resources.SkillPanelSteal:
		return newAreaSteal(objID, arg, core, animMgr)
	case resources.SkillCountBomb:
		return newCountBomb(objID, arg, core, animMgr)
	case resources.SkillTornado:
		return newTornado(objID, arg, core, animMgr)
	case resources.SkillFailed:
		return newFailed(objID, arg, animMgr)
	case resources.SkillQuickGauge:
		return newQuickGauge(objID, arg, core, animMgr)
	case resources.SkillCirkillShot:
		return newCirkillShot(objID, arg, animMgr)
	case resources.SkillShrimpyAttack:
		return newShrimpyAtk(objID, arg, animMgr)
	case resources.SkillBubbleShot, resources.SkillBubbleV, resources.SkillBubbleSide:
		return newBubbleShot(objID, arg, core, animMgr)
	case resources.SkillForteHellsRollingUp, resources.SkillForteHellsRollingDown:
		return newForteHellsRolling(objID, arg, core, animMgr)
	case resources.SkillForteDarkArmBladeType1, resources.SkillForteDarkArmBladeType2:
		return newForteDarkArmBlade(objID, arg, core, skillID, animMgr)
	case resources.SkillForteShootingBuster:
		return newForteShootingBuster(objID, arg, core, animMgr)
	case resources.SkillForteDarknessOverload:
		return newForteDarknessOverload(objID, arg, core, animMgr)
	case resources.SkillChipForteAnother:
		return newChipForteAnother(objID, arg, core, animMgr)
	case resources.SkillDeathMatch1, resources.SkillDeathMatch2, resources.SkillDeathMatch3:
		return newDeathMatch(objID, arg, core, animMgr)
	case resources.SkillPanelReturn:
		return newPanelReturn(objID, arg, core, animMgr)
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

	}
}

func (p *tmpskill) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
*/
