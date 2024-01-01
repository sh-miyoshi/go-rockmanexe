package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
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
	case skillcore.SkillCannon:
		return newCannon(objID, resources.SkillTypeNormalCannon, arg, core)
	case skillcore.SkillHighCannon:
		return newCannon(objID, resources.SkillTypeHighCannon, arg, core)
	case skillcore.SkillMegaCannon:
		return newCannon(objID, resources.SkillTypeMegaCannon, arg, core)
	case skillcore.SkillMiniBomb:
		return newMiniBomb(objID, arg)
	case skillcore.SkillSword:
		return newSword(objID, resources.SkillTypeSword, arg)
	case skillcore.SkillWideSword:
		return newSword(objID, resources.SkillTypeWideSword, arg)
	case skillcore.SkillLongSword:
		return newSword(objID, resources.SkillTypeLongSword, arg)
	case skillcore.SkillShockWave:
		return newShockWave(objID, false, arg)
	case skillcore.SkillRecover:
		return newRecover(objID, arg)
	case skillcore.SkillSpreadGun:
		return newSpreadGun(objID, arg)
	case skillcore.SkillVulcan1:
		return newVulcan(objID, arg)
	case skillcore.SkillPlayerShockWave:
		return newShockWave(objID, true, arg)
	case skillcore.SkillThunderBall:
		return newThunderBall(objID, arg)
	case skillcore.SkillWideShot:
		return newWideShot(objID, arg)
	case skillcore.SkillBoomerang:
		return newBoomerang(objID, arg)
	case skillcore.SkillWaterBomb:
		return newWaterBomb(objID, arg)
	case skillcore.SkillAquamanShot:
		return newAquamanShot(objID, arg)
	case skillcore.SkillAquaman:
		return newAquaman(objID, arg)
	case skillcore.SkillCrackout:
		return newCrack(objID, crackType1, arg)
	case skillcore.SkillDoubleCrack:
		return newCrack(objID, crackType2, arg)
	case skillcore.SkillTripleCrack:
		return newCrack(objID, crackType3, arg)
	case skillcore.SkillBambooLance:
		return newBambooLance(objID, arg)
	case skillcore.SkillDreamSword:
		return newDreamSword(objID, arg)
	case skillcore.SkillInvisible:
		return newInvisible(objID, arg)
	case skillcore.SkillGarooBreath:
		return newGarooBreath(objID, arg)
	case skillcore.SkillFlamePillarTracking:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeTracking)
	case skillcore.SkillFlamePillarRandom:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeRandom)
	case skillcore.SkillHeatShot:
		return newHeatShot(objID, arg, heatShotTypeShot)
	case skillcore.SkillHeatV:
		return newHeatShot(objID, arg, heatShotTypeV)
	case skillcore.SkillHeatSide:
		return newHeatShot(objID, arg, heatShotTypeSide)
	case skillcore.SkillFlamePillarLine:
		return newFlamePillar(objID, arg, resources.SkillFlamePillarTypeLine)
	case skillcore.SkillAreaSteal:
		return newAreaSteal(objID, arg)
	case skillcore.SkillPanelSteal:
		return newPanelSteal(objID, arg)
	case skillcore.SkillCountBomb:
		return newCountBomb(objID, arg)
	case skillcore.SkillTornado:
		return newTornado(objID, arg)
	case skillcore.SkillFailed:
		return newFailed(objID, arg)
	case skillcore.SkillQuickGauge:
		return newQuickGauge(objID, arg)
	case skillcore.SkillCirkillShot:
		return newCirkillShot(objID, arg)
	}

	system.SetError(fmt.Sprintf("Skill %d is not implemented yet", skillID))
	return nil
}

func GetSkillID(chipID int) int {
	switch chipID {
	case chip.IDCannon:
		return skillcore.SkillCannon
	case chip.IDHighCannon:
		return skillcore.SkillHighCannon
	case chip.IDMegaCannon:
		return skillcore.SkillMegaCannon
	case chip.IDSword:
		return skillcore.SkillSword
	case chip.IDWideSword:
		return skillcore.SkillWideSword
	case chip.IDLongSword:
		return skillcore.SkillLongSword
	case chip.IDMiniBomb:
		return skillcore.SkillMiniBomb
	case chip.IDRecover10:
		return skillcore.SkillRecover
	case chip.IDRecover30:
		return skillcore.SkillRecover
	case chip.IDSpreadGun:
		return skillcore.SkillSpreadGun
	case chip.IDVulcan1:
		return skillcore.SkillVulcan1
	case chip.IDShockWave:
		return skillcore.SkillPlayerShockWave
	case chip.IDThunderBall:
		return skillcore.SkillThunderBall
	case chip.IDWideShot:
		return skillcore.SkillWideShot
	case chip.IDBoomerang1:
		return skillcore.SkillBoomerang
	case chip.IDAquaman:
		return skillcore.SkillAquaman
	case chip.IDCrackout:
		return skillcore.SkillCrackout
	case chip.IDDoubleCrack:
		return skillcore.SkillDoubleCrack
	case chip.IDTripleCrack:
		return skillcore.SkillTripleCrack
	case chip.IDBambooLance:
		return skillcore.SkillBambooLance
	case chip.IDDreamSword:
		return skillcore.SkillDreamSword
	case chip.IDInvisible:
		return skillcore.SkillInvisible
	case chip.IDHeatShot:
		return skillcore.SkillHeatShot
	case chip.IDHeatV:
		return skillcore.SkillHeatV
	case chip.IDHeatSide:
		return skillcore.SkillHeatSide
	case chip.IDFlameLine1, chip.IDFlameLine2, chip.IDFlameLine3:
		return skillcore.SkillFlamePillarLine
	case chip.IDAreaSteal:
		return skillcore.SkillAreaSteal
	case chip.IDPanelSteal:
		return skillcore.SkillPanelSteal
	case chip.IDCountBomb:
		return skillcore.SkillCountBomb
	case chip.IDTornado:
		return skillcore.SkillTornado
	case chip.IDAttack10:
		return skillcore.SkillFailed
	case chip.IDQuickGauge:
		return skillcore.SkillQuickGauge
	}

	system.SetError(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
	return 0
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
	Arg Argument

	count int
}

func newTmpSkill(objID string, arg Argument) *tmpskill {
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
