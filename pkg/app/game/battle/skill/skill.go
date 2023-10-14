package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
)

const (
	SkillCannon int = iota
	SkillHighCannon
	SkillMegaCannon
	SkillMiniBomb
	SkillSword
	SkillWideSword
	SkillLongSword
	SkillShockWave
	SkillRecover
	SkillSpreadGun
	SkillVulcan1
	SkillPlayerShockWave
	SkillThunderBall
	SkillWideShot
	SkillBoomerang
	SkillWaterBomb
	SkillAquamanShot
	SkillAquaman
	SkillCrackout
	SkillDoubleCrack
	SkillTripleCrack
	SkillBambooLance
	SkillDreamSword
	SkillInvisible
	SkillGarooBreath
	SkillFlamePillarRandom
	SkillFlamePillarTracking
	SkillHeatShot
	SkillHeatV
	SkillHeatSide
	SkillFlamePillarLine
	SkillAreaSteal
	SkillPanelSteal
)

type SkillAnim interface {
	anim.Anim

	StopByOwner()
}

type Argument struct {
	OwnerID    string
	Power      uint
	TargetType int
}

var (
	imgPick          []int
	imgThunderBall   []int
	imgBoomerang     []int
	imgDreamSword    []int
	imgGarooBreath   []int
	imgFlamePillar   []int
	imgFlameLineBody []int
	imgHeatShotBody  []int
	imgHeatShotAtk   []int
)

func Init() error {
	// TODO: 重複削除
	if err := loadImages(); err != nil {
		return fmt.Errorf("failed to load skill image: %w", err)
	}
	if err := skilldraw.LoadImages(); err != nil {
		return fmt.Errorf("failed to load skill image: %w", err)
	}

	return nil
}

func End() {
	cleanupImages()
	skilldraw.ClearImages()
}

// Get ...
func Get(skillID int, arg Argument) SkillAnim {
	objID := uuid.New().String()

	switch skillID {
	case SkillCannon:
		return newCannon(objID, resources.SkillTypeNormalCannon, arg)
	case SkillHighCannon:
		return newCannon(objID, resources.SkillTypeHighCannon, arg)
	case SkillMegaCannon:
		return newCannon(objID, resources.SkillTypeMegaCannon, arg)
	case SkillMiniBomb:
		return newMiniBomb(objID, arg)
	case SkillSword:
		return newSword(objID, resources.SkillTypeSword, arg)
	case SkillWideSword:
		return newSword(objID, resources.SkillTypeWideSword, arg)
	case SkillLongSword:
		return newSword(objID, resources.SkillTypeLongSword, arg)
	case SkillShockWave:
		return newShockWave(objID, false, arg)
	case SkillRecover:
		return newRecover(objID, arg)
	case SkillSpreadGun:
		return newSpreadGun(objID, arg)
	case SkillVulcan1:
		return newVulcan(objID, arg)
	case SkillPlayerShockWave:
		return newShockWave(objID, true, arg)
	case SkillThunderBall:
		return newThunderBall(objID, arg)
	case SkillWideShot:
		return newWideShot(objID, arg)
	case SkillBoomerang:
		return newBoomerang(objID, arg)
	case SkillWaterBomb:
		return newWaterBomb(objID, arg)
	case SkillAquamanShot:
		return newAquamanShot(objID, arg)
	case SkillAquaman:
		return newAquaman(objID, arg)
	case SkillCrackout:
		return newCrack(objID, crackType1, arg)
	case SkillDoubleCrack:
		return newCrack(objID, crackType2, arg)
	case SkillTripleCrack:
		return newCrack(objID, crackType3, arg)
	case SkillBambooLance:
		return newBambooLance(objID, arg)
	case SkillDreamSword:
		return newDreamSword(objID, arg)
	case SkillInvisible:
		return newInvisible(objID, arg)
	case SkillGarooBreath:
		return newGarooBreath(objID, arg)
	case SkillFlamePillarTracking:
		return newFlamePillar(objID, arg, flamePillarTypeTracking)
	case SkillFlamePillarRandom:
		return newFlamePillar(objID, arg, flamePillarTypeRandom)
	case SkillHeatShot:
		return newHeatShot(objID, arg, heatShotTypeShot)
	case SkillHeatV:
		return newHeatShot(objID, arg, heatShotTypeV)
	case SkillHeatSide:
		return newHeatShot(objID, arg, heatShotTypeSide)
	case SkillFlamePillarLine:
		return newFlamePillar(objID, arg, flamePillarTypeLine)
	case SkillAreaSteal:
		return newAreaSteal(objID, arg)
	case SkillPanelSteal:
		return newPanelSteal(objID, arg)
	}

	common.SetError(fmt.Sprintf("Skill %d is not implemented yet", skillID))
	return nil
}

func GetSkillID(chipID int) int {
	switch chipID {
	case chip.IDCannon:
		return SkillCannon
	case chip.IDHighCannon:
		return SkillHighCannon
	case chip.IDMegaCannon:
		return SkillMegaCannon
	case chip.IDSword:
		return SkillSword
	case chip.IDWideSword:
		return SkillWideSword
	case chip.IDLongSword:
		return SkillLongSword
	case chip.IDMiniBomb:
		return SkillMiniBomb
	case chip.IDRecover10:
		return SkillRecover
	case chip.IDRecover30:
		return SkillRecover
	case chip.IDSpreadGun:
		return SkillSpreadGun
	case chip.IDVulcan1:
		return SkillVulcan1
	case chip.IDShockWave:
		return SkillPlayerShockWave
	case chip.IDThunderBall:
		return SkillThunderBall
	case chip.IDWideShot:
		return SkillWideShot
	case chip.IDBoomerang1:
		return SkillBoomerang
	case chip.IDAquaman:
		return SkillAquaman
	case chip.IDCrackout:
		return SkillCrackout
	case chip.IDDoubleCrack:
		return SkillDoubleCrack
	case chip.IDTripleCrack:
		return SkillTripleCrack
	case chip.IDBambooLance:
		return SkillBambooLance
	case chip.IDDreamSword:
		return SkillDreamSword
	case chip.IDInvisible:
		return SkillInvisible
	case chip.IDHeatShot:
		return SkillHeatShot
	case chip.IDHeatV:
		return SkillHeatV
	case chip.IDHeatSide:
		return SkillHeatSide
	case chip.IDFlameLine1, chip.IDFlameLine2, chip.IDFlameLine3:
		return SkillFlamePillarLine
	case chip.IDAreaSteal:
		return SkillAreaSteal
	case chip.IDPanelSteal:
		return SkillPanelSteal
	}

	common.SetError(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
	return 0
}

/*
Skill template
package skill

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
	pos := localanim.GetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	n := p.count / delay
	if n < len(img) {
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img[n], true)
	}
}

func (p *tmpskill) Process() (bool, error) {
	p.count++

	max := len(img) * delay
	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *tmpskill) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *tmpskill) StopByOwner() {
	localanim.Delete(p.ID)
}
*/
