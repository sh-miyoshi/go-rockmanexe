package skill

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
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
)

const (
	typeNormalCannon int = iota
	typeHighCannon
	typeMegaCannon
)

const (
	typeSword int = iota
	typeWideSword
	typeLongSword
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
	delaySword      = 3
	delayMiniBomb   = 4
	delayShockWave  = 5
	delayRecover    = 1
	delaySpreadGun  = 2
	delayVulcan     = 2
)

type Argument struct {
	OwnerID    string
	Power      uint
	TargetType int
}

var (
	imgCannonAtk     [3][]int32
	imgCannonBody    [3][]int32
	imgSword         [3][]int32
	imgMiniBomb      []int32
	imgShockWave     []int32
	imgRecover       []int32
	imgSpreadGunAtk  []int32
	imgSpreadGunBody []int32
	imgVulcan        []int32
)

type cannon struct {
	Type       int
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type sword struct {
	Type       int
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type miniBomb struct {
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
	OwnerID    string
	Power      uint
	TargetType int

	count int
	x, y  int
}

type recover struct {
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type spreadGun struct {
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type spreadHit struct {
	Power      uint
	TargetType int

	count int
	x, y  int
}

type vulcan struct {
	OwnerID    string
	Power      uint
	TargetType int
	Times      int

	count    int
	imageNo  int
	atkCount int
	hit      bool
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
}

// Get ...
func Get(skillID int, arg Argument) anim.Anim {
	switch skillID {
	case SkillCannon:
		return &cannon{OwnerID: arg.OwnerID, Type: typeNormalCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillHighCannon:
		return &cannon{OwnerID: arg.OwnerID, Type: typeHighCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillMegaCannon:
		return &cannon{OwnerID: arg.OwnerID, Type: typeMegaCannon, Power: arg.Power, TargetType: arg.TargetType}
	case SkillMiniBomb:
		px, py := field.GetPos(arg.OwnerID)
		return &miniBomb{OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, TargetX: px + 3, TargetY: py}
	case SkillSword:
		return &sword{OwnerID: arg.OwnerID, Type: typeSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillWideSword:
		return &sword{OwnerID: arg.OwnerID, Type: typeWideSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillLongSword:
		return &sword{OwnerID: arg.OwnerID, Type: typeLongSword, Power: arg.Power, TargetType: arg.TargetType}
	case SkillShockWave:
		px, py := field.GetPos(arg.OwnerID)
		return &shockWave{OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, x: px, y: py}
	case SkillRecover:
		return &recover{OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillSpreadGun:
		return &spreadGun{OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType}
	case SkillVulcan1:
		return &vulcan{OwnerID: arg.OwnerID, Power: arg.Power, TargetType: arg.TargetType, Times: 3}
	}

	panic(fmt.Sprintf("Skill %d is not implemented yet", skillID))
}

func GetByChip(chipID int, arg Argument) anim.Anim {
	id := -1
	switch chipID {
	case chip.IDCannon:
		id = SkillCannon
	case chip.IDHighCannon:
		id = SkillHighCannon
	case chip.IDMegaCannon:
		id = SkillMegaCannon
	case chip.IDSword:
		id = SkillSword
	case chip.IDWideSword:
		id = SkillWideSword
	case chip.IDLongSword:
		id = SkillLongSword
	case chip.IDMiniBomb:
		id = SkillMiniBomb
	case chip.IDRecover10:
		id = SkillRecover
	case chip.IDRecover30:
		id = SkillRecover
	case chip.IDSpreadGun:
		id = SkillSpreadGun
	case chip.IDVulcan1:
		id = SkillVulcan1
	default:
		panic(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
	}
	return Get(id, arg)
}

func (p *cannon) Draw() {
	px, py := field.GetPos(p.OwnerID)
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
		px, py := field.GetPos(p.OwnerID)
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

func (p *cannon) DamageProc(dm *damage.Damage) {
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeObject,
	}
}

func (p *sword) Draw() {
	px, py := field.GetPos(p.OwnerID)
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

		px, py := field.GetPos(p.OwnerID)

		dm.PosX = px + 1
		dm.PosY = py
		damage.New(dm)

		switch p.Type {
		case typeSword:
			// No more damage area
		case typeWideSword:
			dm.PosY = py - 1
			damage.New(dm)
			dm.PosY = py + 1
			damage.New(dm)
		case typeLongSword:
			dm.PosX = px + 2
			damage.New(dm)
		}
	}

	if p.count > len(imgSword[p.Type])*delaySword {
		return true, nil
	}
	return false, nil
}

func (p *sword) DamageProc(dm *damage.Damage) {
}

func (p *sword) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeObject,
	}
}

func (p *miniBomb) Draw() {
	n := (p.count / delayMiniBomb) % len(imgMiniBomb)
	if n >= 0 {
		vx := p.baseX + int32(p.dx)
		vy := p.baseY + int32(p.dy)
		dxlib.DrawRotaGraph(vx-38, vy-28, 1, 0, imgMiniBomb[n], dxlib.TRUE)
	}
}

func (p *miniBomb) Process() (bool, error) {
	if p.count == 0 {
		// Initialize
		px, py := field.GetPos(p.OwnerID)
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

func (p *miniBomb) DamageProc(dm *damage.Damage) {
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeObject,
	}
}

func (p *shockWave) Draw() {
	n := (p.count / delayShockWave) % len(imgShockWave)
	if n >= 0 {
		vx, vy := battlecommon.ViewPos(p.x, p.y)
		dxlib.DrawRotaGraph(vx, vy, 1, 0, imgShockWave[n], dxlib.TRUE)
	}
}

func (p *shockWave) Process() (bool, error) {
	n := len(imgShockWave) * delayShockWave
	if p.count%(n) == 0 {
		sound.On(sound.SEShockWave)
		// TODO Player Shock Wave
		p.x--
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

func (p *shockWave) DamageProc(dm *damage.Damage) {
}

func (p *shockWave) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeObject,
	}
}

func (p *recover) Draw() {
	n := (p.count / delayRecover) % len(imgRecover)
	if n >= 0 {
		px, py := field.GetPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x, y, 1, 0, imgRecover[n], dxlib.TRUE)
	}
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SERecover)
		px, py := field.GetPos(p.OwnerID)
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

func (p *recover) DamageProc(dm *damage.Damage) {
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
	}
}

func (p *spreadGun) Draw() {
	n := p.count / delaySpreadGun

	// Show body
	if n < len(imgSpreadGunBody) {
		px, py := field.GetPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+50, y-18, 1, 0, imgSpreadGunBody[n], dxlib.TRUE)
	}

	// Show atk
	n = (p.count - 4) / delaySpreadGun
	if n >= 0 && n < len(imgSpreadGunAtk) {
		px, py := field.GetPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+100, y-20, 1, 0, imgSpreadGunAtk[n], dxlib.TRUE)
	}
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		sound.On(sound.SEGun)

		px, py := field.GetPos(p.OwnerID)
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

func (p *spreadGun) DamageProc(dm *damage.Damage) {
}

func (p *spreadGun) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
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

func (p *spreadHit) DamageProc(dm *damage.Damage) {
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
	}
}

func (p *vulcan) Draw() {
	px, py := field.GetPos(p.OwnerID)
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
			px, py := field.GetPos(p.OwnerID)
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

func (p *vulcan) DamageProc(dm *damage.Damage) {
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
	}
}
