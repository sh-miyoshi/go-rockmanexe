package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
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
)

const (
	delayCannonAtk   = 2
	delayCannonBody  = 6
	delaySword       = 3
	delayMiniBomb    = 4
	delayRecover     = 1
	delaySpreadGun   = 2
	delayVulcan      = 2
	delayPick        = 3
	delayThunderBall = 6
	delayWideShot    = 4
	delayBoomerang   = 8
)

type Argument struct {
	OwnerID    string
	Power      uint
	TargetType int
}

var (
	imgCannonAtk     [TypeCannonMax][]int32
	imgCannonBody    [TypeCannonMax][]int32
	imgSword         [TypeSwordMax][]int32
	imgMiniBomb      []int32
	imgShockWave     []int32
	imgRecover       []int32
	imgSpreadGunAtk  []int32
	imgSpreadGunBody []int32
	imgVulcan        []int32
	imgPick          []int32
	imgThunderBall   []int32
	imgWideShotBody  []int32
	imgWideShotBegin []int32
	imgWideShotMove  []int32
	imgBoomerang     []int32
)

func Init() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int32, 24)
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
		imgMiniBomb = append(imgMiniBomb, tmp[i])
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

	return nil
}

func End() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgCannonAtk[i]); j++ {
			dxlib.DeleteGraph(imgCannonAtk[i][j])
		}
		imgCannonAtk[i] = []int32{}
		for j := 0; j < len(imgCannonBody[i]); j++ {
			dxlib.DeleteGraph(imgCannonBody[i][j])
		}
		imgCannonBody[i] = []int32{}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgSword[i]); j++ {
			dxlib.DeleteGraph(imgSword[i][j])
		}
		imgSword[i] = []int32{}
	}
	for i := 0; i < len(imgMiniBomb); i++ {
		dxlib.DeleteGraph(imgMiniBomb[i])
	}
	imgMiniBomb = []int32{}
	for i := 0; i < len(imgShockWave); i++ {
		dxlib.DeleteGraph(imgShockWave[i])
	}
	imgShockWave = []int32{}
	for i := 0; i < len(imgSpreadGunAtk); i++ {
		dxlib.DeleteGraph(imgSpreadGunAtk[i])
	}
	imgSpreadGunAtk = []int32{}
	for i := 0; i < len(imgSpreadGunBody); i++ {
		dxlib.DeleteGraph(imgSpreadGunBody[i])
	}
	imgSpreadGunBody = []int32{}
	for i := 0; i < len(imgVulcan); i++ {
		dxlib.DeleteGraph(imgVulcan[i])
	}
	imgVulcan = []int32{}
	for i := 0; i < len(imgRecover); i++ {
		dxlib.DeleteGraph(imgRecover[i])
	}
	imgRecover = []int32{}
	for i := 0; i < len(imgPick); i++ {
		dxlib.DeleteGraph(imgPick[i])
	}
	imgPick = []int32{}
	for i := 0; i < len(imgThunderBall); i++ {
		dxlib.DeleteGraph(imgThunderBall[i])
	}
	imgThunderBall = []int32{}
	for i := 0; i < len(imgWideShotBody); i++ {
		dxlib.DeleteGraph(imgWideShotBody[i])
	}
	imgWideShotBody = []int32{}
	for i := 0; i < len(imgWideShotBegin); i++ {
		dxlib.DeleteGraph(imgWideShotBegin[i])
	}
	imgWideShotBegin = []int32{}
	for i := 0; i < len(imgWideShotMove); i++ {
		dxlib.DeleteGraph(imgWideShotMove[i])
	}
	imgWideShotMove = []int32{}
	for i := 0; i < len(imgBoomerang); i++ {
		dxlib.DeleteGraph(imgBoomerang[i])
	}
	imgBoomerang = []int32{}
}

// Get ...
func Get(skillID int, arg Argument) anim.Anim {
	objID := uuid.New().String()

	switch skillID {
	case SkillCannon:
		return &cannon{ID: objID, OwnerID: arg.OwnerID, Type: TypeNormalCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillHighCannon:
		return &cannon{ID: objID, OwnerID: arg.OwnerID, Type: TypeHighCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillMegaCannon:
		return &cannon{ID: objID, OwnerID: arg.OwnerID, Type: TypeMegaCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillMiniBomb:
		return newMiniBomb(objID, arg)
	case SkillSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillWideSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeWideSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillLongSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeLongSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillShockWave:
		px, py := objanim.GetObjPos(arg.OwnerID)
		return &shockWave{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: common.DirectLeft, Speed: 5, x: px, y: py}
	case SkillRecover:
		return &recover{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillSpreadGun:
		return &spreadGun{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillVulcan1:
		return &vulcan{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Times: 3}
	case SkillPlayerShockWave:
		px, py := objanim.GetObjPos(arg.OwnerID)
		return &shockWave{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: common.DirectRight, ShowPick: true, Speed: 3, InitWait: 9, x: px, y: py}
	case SkillThunderBall:
		px, py := objanim.GetObjPos(arg.OwnerID)
		x := px + 1
		if arg.TargetType == damage.TargetPlayer {
			x = px - 1
		}

		max := 6 // debug
		return &thunderBall{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, MaxMoveCount: max, x: x, y: py, prevX: px, prevY: py, nextX: x, nextY: py}
	case SkillWideShot:
		px, py := objanim.GetObjPos(arg.OwnerID)
		direct := common.DirectRight
		nextStep := 8
		if arg.TargetType == damage.TargetPlayer {
			direct = common.DirectLeft
			nextStep = 16
		}
		return &wideShot{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: direct, NextStepCount: nextStep, x: px, y: py, state: wideShotStateBegin}
	case SkillBoomerang:
		return newBoomerang(objID, arg)
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

func (p *tmpskill) Draw() {
	px, py := objanim.GetObjPos(p.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	n := p.count / delay
	if n < len(img) {
		dxlib.DrawRotaGraph(x, y, 1, 0, img[n], dxlib.TRUE)
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
