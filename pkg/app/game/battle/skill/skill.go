package skill

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
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
)

const (
	TypeNormalCannon int = iota
	TypeHighCannon
	TypeMegaCannon

	TypeCannonMax
)

const (
	TypeSword int = iota
	TypeWideSword
	TypeLongSword

	TypeSwordMax
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
)

const (
	thunderBallNextStepCount = 80
)

const (
	wideShotStateBegin int = iota
	wideShotStateMove
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
)

type cannon struct {
	ID         string
	Type       int
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type sword struct {
	ID         string
	Type       int
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type miniBomb struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	TargetX    int
	TargetY    int

	count int
	dist  int
	baseX int32
	baseY int32
	dx    int
	dy    int
}

type shockWave struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	Direct     int
	ShowPick   bool
	Speed      int
	InitWait   int

	count    int
	x, y     int
	showWave bool
}

type recover struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type spreadGun struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type spreadHit struct {
	ID         string
	Power      uint
	TargetType int

	count int
	x, y  int
}

type vulcan struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	Times      int

	count    int
	imageNo  int
	atkCount int
	hit      bool
}

type thunderBall struct {
	ID           string
	OwnerID      string
	Power        uint
	TargetType   int
	MaxMoveCount int

	count            int
	x, y             int
	targetX, targetY int
	beforeX, beforeY int
	moveCount        int
	damageID         string
}

type wideShot struct {
	ID            string
	OwnerID       string
	Power         uint
	TargetType    int
	Direct        int
	NextStepCount int

	state    int
	count    int
	x, y     int
	damageID [3]string
}

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
		px, py := anim.GetObjPos(arg.OwnerID)
		return &miniBomb{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, TargetX: px + 3, TargetY: py}
	case SkillSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillWideSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeWideSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillLongSword:
		return &sword{ID: objID, OwnerID: arg.OwnerID, Type: TypeLongSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillShockWave:
		px, py := anim.GetObjPos(arg.OwnerID)
		return &shockWave{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: common.DirectLeft, Speed: 5, x: px, y: py}
	case SkillRecover:
		return &recover{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillSpreadGun:
		return &spreadGun{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillVulcan1:
		return &vulcan{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Times: 3}
	case SkillPlayerShockWave:
		px, py := anim.GetObjPos(arg.OwnerID)
		return &shockWave{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: common.DirectRight, ShowPick: true, Speed: 3, InitWait: 9, x: px, y: py}
	case SkillThunderBall:
		px, py := anim.GetObjPos(arg.OwnerID)
		tx := px + 1
		if arg.TargetType == damage.TargetPlayer {
			tx = px - 1
		}

		max := 6 // debug
		return &thunderBall{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, MaxMoveCount: max, x: px, y: py, targetX: tx, targetY: py}
	case SkillWideShot:
		px, py := anim.GetObjPos(arg.OwnerID)
		direct := common.DirectRight
		nextStep := 8
		if arg.TargetType == damage.TargetPlayer {
			direct = common.DirectLeft
			nextStep = 16
		}
		return &wideShot{ID: objID, OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Direct: direct, NextStepCount: nextStep, x: px, y: py, state: wideShotStateBegin}
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
	}

	panic(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
}

func (p *cannon) Draw() {
	px, py := anim.GetObjPos(p.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	n := p.count / delayCannonBody
	if n < len(imgCannonBody[p.Type]) {
		if n >= 3 {
			x -= 15
		}

		dxlib.DrawRotaGraph(x+48, y-12, 1, 0, imgCannonBody[p.Type][n], dxlib.TRUE)
	}

	n = (p.count - 15) / delayCannonAtk
	if n >= 0 && n < len(imgCannonAtk[p.Type]) {
		dxlib.DrawRotaGraph(x+90, y-10, 1, 0, imgCannonAtk[p.Type][n], dxlib.TRUE)
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		sound.On(sound.SECannon)
		px, py := anim.GetObjPos(p.OwnerID)
		dm := damage.Damage{
			PosY:          py,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeCannonHit,
		}

		if p.TargetType == damage.TargetEnemy {
			for x := px + 1; x < field.FieldNumX; x++ {
				dm.PosX = x
				if field.GetPanelInfo(x, dm.PosY).ObjectID != "" {
					damage.New(dm)
					break
				}
			}
		} else {
			for x := px - 1; x >= 0; x-- {
				dm.PosX = x
				if field.GetPanelInfo(x, dm.PosY).ObjectID != "" {
					damage.New(dm)
					break
				}
			}
		}
	}

	max := len(imgCannonBody[p.Type]) * delayCannonBody
	if max < len(imgCannonAtk[p.Type])*delayCannonAtk+15 {
		max = len(imgCannonAtk[p.Type])*delayCannonAtk + 15
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *sword) Draw() {
	px, py := anim.GetObjPos(p.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	n := (p.count - 5) / delaySword
	if n >= 0 && n < len(imgSword[p.Type]) {
		dxlib.DrawRotaGraph(x+100, y, 1, 0, imgSword[p.Type][n], dxlib.TRUE)
	}
}

func (p *sword) Process() (bool, error) {
	p.count++

	if p.count == 1*delaySword {
		sound.On(sound.SESword)

		dm := damage.Damage{
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		}

		px, py := anim.GetObjPos(p.OwnerID)

		dm.PosX = px + 1
		dm.PosY = py
		damage.New(dm)

		switch p.Type {
		case TypeSword:
			// No more damage area
		case TypeWideSword:
			dm.PosY = py - 1
			damage.New(dm)
			dm.PosY = py + 1
			damage.New(dm)
		case TypeLongSword:
			dm.PosX = px + 2
			damage.New(dm)
		}
	}

	if p.count > len(imgSword[p.Type])*delaySword {
		return true, nil
	}
	return false, nil
}

func (p *sword) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *sword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *miniBomb) Draw() {
	n := (p.count / delayMiniBomb) % len(imgMiniBomb)
	if n >= 0 {
		vx := p.baseX + int32(p.dx)
		vy := p.baseY + int32(p.dy)
		dxlib.DrawRotaGraph(vx-38, vy, 1, 0, imgMiniBomb[n], dxlib.TRUE)
	}
}

func (p *miniBomb) Process() (bool, error) {
	if p.count == 0 {
		// Initialize
		px, py := anim.GetObjPos(p.OwnerID)
		p.baseX, p.baseY = battlecommon.ViewPos(px, py)
		// TODO: yが等しい場合でかつプレイヤー側のみ
		p.dist = (p.TargetX - px) * field.PanelSizeX

		sound.On(sound.SEBombThrow)
	}

	// y = ax^2 + bx +c
	// (0,0), (d/2, ymax), (d, 0)
	p.count++
	p.dx += 4
	ymax := 100
	p.dy = ymax*4*p.dx*p.dx/(p.dist*p.dist) - ymax*4*p.dx/p.dist

	if p.dx >= p.dist+38 {
		// TODO 不発処理(画面外やパネル状況など)
		sound.On(sound.SEExplode)
		anim.New(effect.Get(effect.TypeExplode, p.TargetX, p.TargetY, 0))
		damage.New(damage.Damage{
			PosX:          p.TargetX,
			PosY:          p.TargetY,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		})
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *shockWave) Draw() {
	n := (p.count / p.Speed) % len(imgShockWave)
	if p.showWave && n >= 0 {
		vx, vy := battlecommon.ViewPos(p.x, p.y)
		if p.Direct == common.DirectLeft {
			dxlib.DrawRotaGraph(vx, vy, 1, 0, imgShockWave[n], dxlib.TRUE)
		} else if p.Direct == common.DirectRight {
			ofsx := int32(100 / 2)
			ofsy := int32(140 / 2)
			dxlib.DrawTurnGraph(vx-ofsx, vy-ofsy, imgShockWave[n], dxlib.TRUE)
		}
	}

	if p.ShowPick {
		n = (p.count / delayPick)
		if n < len(imgPick) {
			px, py := anim.GetObjPos(p.OwnerID)
			vx, vy := battlecommon.ViewPos(px, py)
			dxlib.DrawRotaGraph(vx, vy-15, 1, 0, imgPick[n], dxlib.TRUE)
		}
	}
}

func (p *shockWave) Process() (bool, error) {
	if p.count < p.InitWait {
		p.count++
		return false, nil
	}

	n := len(imgShockWave) * p.Speed
	if p.count%(n) == 0 {
		p.showWave = true
		sound.On(sound.SEShockWave)
		if p.Direct == common.DirectLeft {
			p.x--
		} else if p.Direct == common.DirectRight {
			p.x++
		}
		damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           n - 2,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
		})
	}
	p.count++

	if p.x < 0 || p.x > field.FieldNumX {
		return true, nil
	}
	return false, nil
}

func (p *shockWave) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *shockWave) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *recover) Draw() {
	n := (p.count / delayRecover) % len(imgRecover)
	if n >= 0 {
		px, py := anim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x, y, 1, 0, imgRecover[n], dxlib.TRUE)
	}
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SERecover)
		px, py := anim.GetObjPos(p.OwnerID)
		damage.New(damage.Damage{
			PosX:          px,
			PosY:          py,
			Power:         -int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		})
	}

	p.count++

	if p.count > len(imgRecover)*delayRecover {
		return true, nil
	}
	return false, nil
}

func (p *recover) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *spreadGun) Draw() {
	n := p.count / delaySpreadGun

	// Show body
	if n < len(imgSpreadGunBody) {
		px, py := anim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+50, y-18, 1, 0, imgSpreadGunBody[n], dxlib.TRUE)
	}

	// Show atk
	n = (p.count - 4) / delaySpreadGun
	if n >= 0 && n < len(imgSpreadGunAtk) {
		px, py := anim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+100, y-20, 1, 0, imgSpreadGunAtk[n], dxlib.TRUE)
	}
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		sound.On(sound.SEGun)

		px, py := anim.GetObjPos(p.OwnerID)
		for x := px + 1; x < field.FieldNumX; x++ {
			if field.GetPanelInfo(x, py).ObjectID != "" {
				// Hit
				sound.On(sound.SESpreadHit)

				damage.New(damage.Damage{
					PosX:          x,
					PosY:          py,
					Power:         int(p.Power),
					TTL:           1,
					TargetType:    p.TargetType,
					HitEffectType: effect.TypeHitBig,
				})
				// Spreading
				for sy := -1; sy <= 1; sy++ {
					if py+sy < 0 || py+sy >= field.FieldNumY {
						continue
					}
					for sx := -1; sx <= 1; sx++ {
						if sy == 0 && sx == 0 {
							continue
						}
						if x+sx >= 0 && x+sx < field.FieldNumX {
							anim.New(&spreadHit{
								Power:      p.Power,
								TargetType: p.TargetType,
								x:          x + sx,
								y:          py + sy,
							})
						}
					}
				}

				break
			}
		}
	}

	p.count++

	max := len(imgSpreadGunAtk)
	if len(imgSpreadGunBody) > max {
		max = len(imgSpreadGunBody)
	}

	if p.count > max*delaySpreadGun {
		return true, nil
	}
	return false, nil
}

func (p *spreadGun) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *spreadGun) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *spreadHit) Draw() {
}

func (p *spreadHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		anim.New(effect.Get(effect.TypeSpreadHit, p.x, p.y, 5))
		damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		})

		return true, nil
	}
	return false, nil
}

func (p *spreadHit) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *vulcan) Draw() {
	px, py := anim.GetObjPos(p.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	// Show body
	dxlib.DrawRotaGraph(x+50, y-18, 1, 0, imgVulcan[p.imageNo], dxlib.TRUE)
	// Show attack
	if p.imageNo != 0 {
		if p.imageNo%2 == 0 {
			dxlib.DrawRotaGraph(x+100, y-10, 1, 0, imgVulcan[3], dxlib.TRUE)
		} else {
			dxlib.DrawRotaGraph(x+100, y-15, 1, 0, imgVulcan[3], dxlib.TRUE)
		}
	}
}

func (p *vulcan) Process() (bool, error) {
	p.count++
	if p.count >= delayVulcan*1 {
		if p.count%(delayVulcan*5) == delayVulcan*1 {
			sound.On(sound.SEGun)

			p.imageNo = p.imageNo%2 + 1
			// Add damage
			px, py := anim.GetObjPos(p.OwnerID)
			hit := false
			for x := px + 1; x < field.FieldNumX; x++ {
				if field.GetPanelInfo(x, py).ObjectID != "" {
					damage.New(damage.Damage{
						PosX:          x,
						PosY:          py,
						Power:         int(p.Power),
						TTL:           1,
						TargetType:    p.TargetType,
						HitEffectType: effect.TypeSpreadHit,
					})
					anim.New(effect.Get(effect.TypeVulcanHit1, x, py, 20))
					if p.hit && x < field.FieldNumX-1 {
						anim.New(effect.Get(effect.TypeVulcanHit2, x+1, py, 20))
						damage.New(damage.Damage{
							PosX:          x + 1,
							PosY:          py,
							Power:         int(p.Power),
							TTL:           1,
							TargetType:    p.TargetType,
							HitEffectType: effect.TypeNone,
						})
					}
					hit = true
					sound.On(sound.SECannonHit)
					break
				}
			}
			p.hit = hit
			p.atkCount++
			if p.atkCount == p.Times {
				return true, nil
			}
		}

	}

	return false, nil
}

func (p *vulcan) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *thunderBall) Draw() {
	x, y := battlecommon.ViewPos(p.x, p.y)
	n := (p.count / delayThunderBall) % len(imgThunderBall)

	c := p.count % thunderBallNextStepCount
	if c != 0 {
		ofsx := battlecommon.GetOffset(p.targetX, p.x, p.beforeX, c, thunderBallNextStepCount, field.PanelSizeX)
		ofsy := battlecommon.GetOffset(p.targetY, p.y, p.beforeY, c, thunderBallNextStepCount, field.PanelSizeY)
		dxlib.DrawRotaGraph(x+int32(ofsx), y+25+int32(ofsy), 1, 0, imgThunderBall[n], dxlib.TRUE)
	}
}

func (p *thunderBall) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SEThunderBall)
	}

	halfNext := thunderBallNextStepCount / 2
	if p.damageID != "" {
		if !damage.Exists(p.damageID) && p.count%halfNext != 0 {
			// attack hit to target
			return true, nil
		}
	}

	if p.count%thunderBallNextStepCount == 0 {
		// Set current position
		p.beforeX = p.x
		p.beforeY = p.y
		p.x = p.targetX
		p.y = p.targetY

		// Decide next position
		objType := anim.ObjTypePlayer
		if p.TargetType == damage.TargetEnemy {
			objType = anim.ObjTypeEnemy
		}

		objs := anim.GetObjs(anim.Filter{ObjType: objType})
		if len(objs) == 0 {
			// no target
			if p.TargetType == damage.TargetPlayer {
				p.targetX--
			} else {
				p.targetX++
			}
		} else {
			xdif := objs[0].PosX - p.x
			ydif := objs[0].PosY - p.y

			if xdif != 0 || ydif != 0 {
				if common.Abs(xdif) > common.Abs(ydif) {
					// move to x
					p.targetX = p.x + (xdif / common.Abs(xdif))
					p.targetY = p.y
				} else {
					// move to y
					p.targetX = p.x
					p.targetY = p.y + (ydif / common.Abs(ydif))
				}
			}
		}

		p.moveCount++
		if p.moveCount > p.MaxMoveCount {
			return true, nil
		}

		if p.x < 0 || p.x > field.FieldNumX || p.y < 0 || p.y > field.FieldNumY {
			return true, nil
		}
	}
	if p.count%halfNext == 0 {
		p.damageID = damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           halfNext,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
		})
	}

	p.count++
	return false, nil
}

func (p *thunderBall) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *thunderBall) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}

func (p *wideShot) Draw() {
	opt := dxlib.DrawRotaGraphOption{}
	ofs := int32(1)
	if p.Direct == common.DirectLeft {
		xflip := int32(dxlib.TRUE)
		opt.ReverseXFlag = &xflip
		ofs = -1
	}

	switch p.state {
	case wideShotStateBegin:
		x, y := battlecommon.ViewPos(p.x, p.y)
		n := (p.count / delayWideShot)

		if n < len(imgWideShotBody) && p.TargetType == damage.TargetEnemy {
			dxlib.DrawRotaGraph(x+40, y-13, 1, 0, imgWideShotBody[n], dxlib.TRUE, opt)
		}
		if n >= len(imgWideShotBegin) {
			n = len(imgWideShotBegin) - 1
		}
		dxlib.DrawRotaGraph(x+62*ofs, y+20, 1, 0, imgWideShotBegin[n], dxlib.TRUE, opt)
	case wideShotStateMove:
		x, y := battlecommon.ViewPos(p.x, p.y)
		n := (p.count / delayWideShot) % len(imgWideShotMove)
		next := p.x + 1
		prev := p.x - 1
		if p.Direct == common.DirectLeft {
			next, prev = prev, next
		}

		c := p.count % p.NextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, p.x, prev, c, p.NextStepCount, field.PanelSizeX)
			dxlib.DrawRotaGraph(x+int32(ofsx), y+20, 1, 0, imgWideShotMove[n], dxlib.TRUE, opt)
		}
	}
}

func (p *wideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !damage.Exists(did) && p.count%p.NextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.state {
	case wideShotStateBegin:
		if p.count == 0 {
			sound.On(sound.SEWideShot)
		}

		max := len(imgWideShotBody)
		if len(imgWideShotBegin) > max {
			max = len(imgWideShotBegin)
		}
		max *= delayWideShot
		if p.count > max {
			p.state = wideShotStateMove
			p.count = 0
			return false, nil
		}
	case wideShotStateMove:
		if p.count%p.NextStepCount == 0 {
			if p.Direct == common.DirectRight {
				p.x++
			} else if p.Direct == common.DirectLeft {
				p.x--
			}

			if p.x >= field.FieldNumX || p.x < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.y + i
				if y < 0 || y >= field.FieldNumY {
					continue
				}

				p.damageID[i+1] = damage.New(damage.Damage{
					PosX:          p.x,
					PosY:          y,
					Power:         int(p.Power),
					TTL:           p.NextStepCount,
					TargetType:    p.TargetType,
					HitEffectType: effect.TypeNone,
				})
			}
		}
	}

	p.count++
	return false, nil
}

func (p *wideShot) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *wideShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.TypeSkill,
		ObjType:  anim.ObjTypeNone,
	}
}
