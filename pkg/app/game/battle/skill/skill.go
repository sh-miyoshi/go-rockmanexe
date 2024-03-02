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
	arg.PanelBreak = field.PanelBreak
	core := localanim.SkillManager().Get(skillID, arg)

	switch skillID {
	case resources.SkillCannon, resources.SkillHighCannon, resources.SkillMegaCannon:
		return newCannon(objID, arg, core)
	case resources.SkillMiniBomb:
		return newMiniBomb(objID, arg, core)
	case resources.SkillSword, resources.SkillWideSword, resources.SkillLongSword, resources.SkillDreamSword:
		return newSword(objID, arg, core)
	case resources.SkillPlayerShockWave, resources.SkillEnemyShockWave:
		return newShockWave(objID, arg, core, skillID == resources.SkillPlayerShockWave)
	case resources.SkillRecover:
		return newRecover(objID, arg, core)
	case resources.SkillSpreadGun:
		return newSpreadGun(objID, arg, core)
	case resources.SkillVulcan1:
		return newVulcan(objID, arg, core)
	case resources.SkillThunderBall:
		return newThunderBall(objID, arg)
	case resources.SkillPlayerWideShot, resources.SkillEnemyWideShot:
		return newWideShot(objID, arg, core, skillID == resources.SkillPlayerWideShot)
	case resources.SkillBoomerang:
		return newBoomerang(objID, arg, core)
	case resources.SkillWaterBomb:
		return newWaterBomb(objID, arg, core)
	case resources.SkillAquamanShot:
		return newAquamanShot(objID, arg)
	case resources.SkillAquaman:
		return newAquaman(objID, arg)
	case resources.SkillCrackout, resources.SkillDoubleCrack, resources.SkillTripleCrack:
		return newCrack(objID, arg, core)
	case resources.SkillBambooLance:
		return newBambooLance(objID, arg, core)
	case resources.SkillInvisible:
		return newInvisible(objID, arg)
	case resources.SkillGarooBreath:
		return newGarooBreath(objID, arg)
	case resources.SkillFlamePillarTracking, resources.SkillFlamePillarRandom, resources.SkillFlamePillarLine:
		return newFlamePillar(objID, arg, core)
	case resources.SkillHeatShot, resources.SkillHeatV, resources.SkillHeatSide:
		return newHeatShot(objID, arg, core)
	case resources.SkillAreaSteal:
		return newAreaSteal(objID, arg)
	case resources.SkillPanelSteal:
		return newPanelSteal(objID, arg)
	case resources.SkillCountBomb:
		return newCountBomb(objID, arg)
	case resources.SkillTornado:
		return newTornado(objID, arg, core)
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
