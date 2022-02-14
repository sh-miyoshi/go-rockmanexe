package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
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
)

const (
	delayCannonAtk   = 2
	delayCannonBody  = 6
	delaySword       = 3
	delayBombThrow   = 4
	delayRecover     = 1
	delaySpreadGun   = 2
	delayVulcan      = 2
	delayPick        = 3
	delayThunderBall = 6
	delayWideShot    = 4
	delayBoomerang   = 8
	delayBambooLance = 4
)

type Argument struct {
	OwnerID    string
	Power      uint
	TargetType int
}

var (
	imgCannonAtk     [TypeCannonMax][]int
	imgCannonBody    [TypeCannonMax][]int
	imgSword         [TypeSwordMax][]int
	imgBombThrow     []int
	imgShockWave     []int
	imgRecover       []int
	imgSpreadGunAtk  []int
	imgSpreadGunBody []int
	imgVulcan        []int
	imgPick          []int
	imgThunderBall   []int
	imgWideShotBody  []int
	imgWideShotBegin []int
	imgWideShotMove  []int
	imgBoomerang     []int
	imgAquamanShot   []int
	imgBambooLance   []int
	imgDreamSword    []int
)

func Init() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int, 24)
	fname := path + "キャノン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCannonAtk[0] = append(imgCannonAtk[0], tmp[i])
		imgCannonAtk[1] = append(imgCannonAtk[1], tmp[i+8])
		imgCannonAtk[2] = append(imgCannonAtk[2], tmp[i+16])
	}
	fname = path + "キャノン_body.png"
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgCannonBody[0] = append(imgCannonBody[0], tmp[i])
		imgCannonBody[1] = append(imgCannonBody[1], tmp[i+5])
		imgCannonBody[2] = append(imgCannonBody[2], tmp[i+10])
	}

	fname = path + "ミニボム.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgBombThrow = append(imgBombThrow, tmp[i])
	}

	fname = path + "ソード.png"
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		// Note: In the image, the order of wide sword and long sword is swapped.
		imgSword[0] = append(imgSword[0], tmp[i])
		imgSword[1] = append(imgSword[1], tmp[i+8])
		imgSword[2] = append(imgSword[2], tmp[i+4])
	}

	fname = path + "ショックウェーブ.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 7; i++ {
		imgShockWave = append(imgShockWave, tmp[i])
	}

	fname = path + "リカバリー.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgRecover = append(imgRecover, tmp[i])
	}

	fname = path + "スプレッドガン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgSpreadGunAtk = append(imgSpreadGunAtk, tmp[i])
	}
	fname = path + "スプレッドガン_body.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgSpreadGunBody = append(imgSpreadGunBody, tmp[i])
	}

	fname = path + "バルカン.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgVulcan = append(imgVulcan, tmp[i])
	}

	fname = path + "ウェーブ_body.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgPick = append(imgPick, tmp[i])
	}
	imgPick = append(imgPick, tmp[3])
	imgPick = append(imgPick, tmp[3])

	fname = path + "サンダーボール.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgThunderBall = append(imgThunderBall, tmp[i])
	}

	fname = path + "ワイドショット_body.png"
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 3; i++ {
		imgWideShotBody = append(imgWideShotBody, tmp[i])
	}
	fname = path + "ワイドショット_begin.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgWideShotBegin = append(imgWideShotBegin, tmp[i])
	}
	fname = path + "ワイドショット_move.png"
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 3; i++ {
		imgWideShotMove = append(imgWideShotMove, tmp[i])
	}

	fname = path + "ブーメラン.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 80, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgBoomerang = append(imgBoomerang, tmp[i])
	}
	fname = path + "aquaman_shot.png"
	imgAquamanShot = make([]int, 1)
	if imgAquamanShot[0] = dxlib.LoadGraph(fname); imgAquamanShot[0] == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "バンブーランス.png"
	imgBambooLance = make([]int, 1)
	if imgBambooLance[0] = dxlib.LoadGraph(fname); imgBambooLance[0] == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "ドリームソード.png"
	imgDreamSword = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 200, 188, imgDreamSword); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func End() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgCannonAtk[i]); j++ {
			dxlib.DeleteGraph(imgCannonAtk[i][j])
		}
		imgCannonAtk[i] = []int{}
		for j := 0; j < len(imgCannonBody[i]); j++ {
			dxlib.DeleteGraph(imgCannonBody[i][j])
		}
		imgCannonBody[i] = []int{}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgSword[i]); j++ {
			dxlib.DeleteGraph(imgSword[i][j])
		}
		imgSword[i] = []int{}
	}
	for i := 0; i < len(imgBombThrow); i++ {
		dxlib.DeleteGraph(imgBombThrow[i])
	}
	imgBombThrow = []int{}
	for i := 0; i < len(imgShockWave); i++ {
		dxlib.DeleteGraph(imgShockWave[i])
	}
	imgShockWave = []int{}
	for i := 0; i < len(imgSpreadGunAtk); i++ {
		dxlib.DeleteGraph(imgSpreadGunAtk[i])
	}
	imgSpreadGunAtk = []int{}
	for i := 0; i < len(imgSpreadGunBody); i++ {
		dxlib.DeleteGraph(imgSpreadGunBody[i])
	}
	imgSpreadGunBody = []int{}
	for i := 0; i < len(imgVulcan); i++ {
		dxlib.DeleteGraph(imgVulcan[i])
	}
	imgVulcan = []int{}
	for i := 0; i < len(imgRecover); i++ {
		dxlib.DeleteGraph(imgRecover[i])
	}
	imgRecover = []int{}
	for i := 0; i < len(imgPick); i++ {
		dxlib.DeleteGraph(imgPick[i])
	}
	imgPick = []int{}
	for i := 0; i < len(imgThunderBall); i++ {
		dxlib.DeleteGraph(imgThunderBall[i])
	}
	imgThunderBall = []int{}
	for i := 0; i < len(imgWideShotBody); i++ {
		dxlib.DeleteGraph(imgWideShotBody[i])
	}
	imgWideShotBody = []int{}
	for i := 0; i < len(imgWideShotBegin); i++ {
		dxlib.DeleteGraph(imgWideShotBegin[i])
	}
	imgWideShotBegin = []int{}
	for i := 0; i < len(imgWideShotMove); i++ {
		dxlib.DeleteGraph(imgWideShotMove[i])
	}
	imgWideShotMove = []int{}
	for i := 0; i < len(imgBoomerang); i++ {
		dxlib.DeleteGraph(imgBoomerang[i])
	}
	imgBoomerang = []int{}
	for i := 0; i < len(imgAquamanShot); i++ {
		dxlib.DeleteGraph(imgAquamanShot[i])
	}
	imgAquamanShot = []int{}
	for i := 0; i < len(imgBambooLance); i++ {
		dxlib.DeleteGraph(imgBambooLance[i])
	}
	imgBambooLance = []int{}
	for i := 0; i < len(imgDreamSword); i++ {
		dxlib.DeleteGraph(imgDreamSword[i])
	}
	imgDreamSword = []int{}
}

// Get ...
func Get(skillID int, arg Argument) anim.Anim {
	objID := uuid.New().String()

	switch skillID {
	case SkillCannon:
		return newCannon(objID, TypeNormalCannon, arg)
	case SkillHighCannon:
		return newCannon(objID, TypeHighCannon, arg)
	case SkillMegaCannon:
		return newCannon(objID, TypeMegaCannon, arg)
	case SkillMiniBomb:
		return newMiniBomb(objID, arg)
	case SkillSword:
		return newSword(objID, TypeSword, arg)
	case SkillWideSword:
		return newSword(objID, TypeWideSword, arg)
	case SkillLongSword:
		return newSword(objID, TypeLongSword, arg)
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
		res, err := newAquaman(objID, arg)
		if err != nil {
			panic(err)
		}
		return res
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
	}

	panic(fmt.Sprintf("Skill %d is not implemented yet", skillID))
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
	}

	panic(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
}

/*
Skill template
package skill

type tmpskill struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func newTmpSkill(objID string, arg Argument) *tmpskill {
	return &tmpskill{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
	}
}

func (p *tmpskill) Draw() {
	pos := objanim.GetObjPos(p.OwnerID)
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
		AnimType: anim.AnimTypeSkill,
	}
}

*/
