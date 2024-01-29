package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
)

type SkillAnim interface {
	anim.Anim

	StopByOwner()
}

func Init() error {
	if err := skilldraw.LoadImages(); err != nil {
		return fmt.Errorf("failed to load skill image: %w", err)
	}

	return nil
}

func End() {
	skilldraw.ClearImages()
}

func Get(skillID int, arg skillcore.Argument) SkillAnim {
	objID := uuid.New().String()
	arg.GetPanelInfo = field.GetPanelInfo
	core := localanim.SkillManager().Get(skillID, arg)

	switch skillID {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(objID, skillID, arg, core)
	case resources.SkillMiniBomb:
		return newMiniBomb(objID, arg, core)
	case resources.SkillSword:
		return newSword(objID, resources.SkillTypeSword, arg)
	case resources.SkillWideSword:
		return newSword(objID, resources.SkillTypeWideSword, arg)
	case resources.SkillLongSword:
		return newSword(objID, resources.SkillTypeLongSword, arg)
	case resources.SkillEnemyShockWave:
		return newShockWave(objID, false, arg, core)
	case resources.SkillRecover:
		return newRecover(objID, arg, core)
	case resources.SkillSpreadGun:
		return newSpreadGun(objID, arg)
	case resources.SkillVulcan1:
		return newVulcan(objID, arg)
	case resources.SkillPlayerShockWave:
		return newShockWave(objID, true, arg, core)
	case resources.SkillThunderBall:
		return newThunderBall(objID, arg)
	case resources.SkillWideShot:
		return newWideShot(objID, arg)
	case resources.SkillBoomerang:
		return newBoomerang(objID, arg)
	case resources.SkillWaterBomb:
		return newWaterBomb(objID, arg)
	case resources.SkillAquamanShot:
		return newAquamanShot(objID, arg)
	case resources.SkillAquaman:
		return newAquaman(objID, arg)
	case resources.SkillCrackout:
		return newCrack(objID, crackType1, arg)
	case resources.SkillDoubleCrack:
		return newCrack(objID, crackType2, arg)
	case resources.SkillTripleCrack:
		return newCrack(objID, crackType3, arg)
	case resources.SkillBambooLance:
		return newBambooLance(objID, arg)
	case resources.SkillDreamSword:
		return newDreamSword(objID, arg)
	case resources.SkillInvisible:
		return newInvisible(objID, arg)
	case resources.SkillGarooBreath:
		return newGarooBreath(objID, arg)
	case resources.SkillFlamePillarTracking:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeTracking)
	case resources.SkillFlamePillarRandom:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeRandom)
	case resources.SkillHeatShot:
		return newHeatShot(objID, arg, heatShotTypeShot)
	case resources.SkillHeatV:
		return newHeatShot(objID, arg, heatShotTypeV)
	case resources.SkillHeatSide:
		return newHeatShot(objID, arg, heatShotTypeSide)
	case resources.SkillFlamePillarLine:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeLine)
	case resources.SkillAreaSteal:
		return newAreaSteal(objID, arg)
	case resources.SkillPanelSteal:
		return newPanelSteal(objID, arg)
	case resources.SkillCountBomb:
		return newCountBomb(objID, arg)
	case resources.SkillTornado:
		return newTornado(objID, arg)
	case resources.SkillFailed:
		return newFailed(objID, arg)
	case resources.SkillQuickGauge:
		return newQuickGauge(objID, arg)
	case resources.SkillCirkillShot:
		return newCirkillShot(objID, arg)
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
)

type tmpskill struct {
	ID  string
	Arg skillcore.Argument

	count int
}

func newTmpSkill(objID string, arg skillcore.Argument) *tmpskill {
	return &tmpskill{
		ID:  objID,
		Arg: arg,
	}
}

func (p *tmpskill) Draw() {
	// p.drawer.Draw()
}

func (p *tmpskill) Process() (bool, error) {
	p.count++

	return false, nil
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
