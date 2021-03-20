package skill

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	// Chip Base Skills

	SkillCannon int = iota
	SkillHighCannon
	SkillMegaCannon
	SkillMiniBomb
	SkillSword
	SkillWideSword
	SkillLongSword

	skillMax
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
	delayCannonBody = 5
	delaySword      = 3
	delayMiniBomb   = 6
)

type Argument struct {
	OwnerID    string
	Power      int
	TargetType int
	TargetX    int
	TargetY    int
}

var (
	imgCannonAtk  [3][]int32
	imgCannonBody [3][]int32
	imgSword      [3][]int32
	imgMiniBomb   []int32
)

type cannon struct {
	Type       int
	OwnerID    string
	Power      int
	TargetType int

	count int
}

type sword struct {
	Type       int
	OwnerID    string
	Power      int
	TargetType int

	count int
}

type miniBomb struct {
	OwnerID    string
	Power      int
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

func Init() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int32, 24)
	fname := path + "キャノン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCannonAtk[0] = append(imgCannonAtk[0], tmp[i])
		imgCannonAtk[1] = append(imgCannonAtk[1], tmp[i+8])
		imgCannonAtk[2] = append(imgCannonAtk[2], tmp[i+16])
	}
	fname = path + "キャノン_body.png"
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgCannonBody[0] = append(imgCannonBody[0], tmp[i])
		imgCannonBody[1] = append(imgCannonBody[1], tmp[i+5])
		imgCannonBody[2] = append(imgCannonBody[2], tmp[i+10])
	}

	fname = path + "ミニボム.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, tmp); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgMiniBomb = append(imgMiniBomb, tmp[i])
	}

	fname = path + "ソード.png"
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		// Note: In the image, the order of wide sword and long sword is swapped.
		imgSword[0] = append(imgSword[0], tmp[i])
		imgSword[1] = append(imgSword[1], tmp[i+8])
		imgSword[2] = append(imgSword[2], tmp[i+4])
	}

	return nil
}

func End() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgCannonAtk[i]); j++ {
			dxlib.DeleteGraph(imgCannonAtk[i][j])
		}
		for j := 0; j < len(imgCannonBody[i]); j++ {
			dxlib.DeleteGraph(imgCannonBody[i][j])
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgSword[i]); j++ {
			dxlib.DeleteGraph(imgSword[i][j])
		}
	}
	for i := 0; i < len(imgMiniBomb); i++ {
		dxlib.DeleteGraph(imgMiniBomb[i])
	}
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
	}

	panic(fmt.Sprintf("Skill %d is not implemented yet", skillID))
}

func (p *cannon) Draw() {
	px, py := field.GetPos(p.OwnerID)
	if px < 0 || py < 0 {
		logger.Error("Failed to get object %s position", p.OwnerID)
		return
	}
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
		dxlib.DrawRotaGraph(x+115, y-10, 1, 0, imgCannonAtk[p.Type][n], dxlib.TRUE)
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++
	if p.count == 20 {
		px, py := field.GetPos(p.OwnerID)
		dm := damage.Damage{
			PosY:          py,
			Power:         p.Power,
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeHitBig, // TODO
		}

		if p.TargetType == damage.TargetEnemy {
			for x := px + 1; x < field.FieldNumX; x++ {
				dm.PosX = x
				damage.New(dm)
			}
		} else {
			for x := px - 1; x >= 0; x-- {
				dm.PosX = x
				damage.New(dm)
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

func (p *sword) Draw() {
	px, py := field.GetPos(p.OwnerID)
	if px < 0 || py < 0 {
		logger.Error("Failed to get object %s position", p.OwnerID)
		return
	}
	x, y := battlecommon.ViewPos(px, py)

	n := (p.count - 5) / delaySword
	if n >= 0 && n < len(imgSword[p.Type]) {
		dxlib.DrawRotaGraph(x+100, y, 1, 0, imgSword[p.Type][n], dxlib.TRUE)
	}
}

func (p *sword) Process() (bool, error) {
	p.count++

	// TODO damage register

	if p.count > len(imgSword[p.Type])*delaySword {
		return true, nil
	}
	return false, nil
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
	}

	// y = ax^2 + bx +c
	// (0,0), (d/2, ymax), (d, 0)
	p.count++
	p.dx += 4
	ymax := 100
	p.dy = ymax*4*p.dx*p.dx/(p.dist*p.dist) - ymax*4*p.dx/p.dist

	if p.dx >= p.dist+38 {
		// TODO damage register
		return true, nil
	}
	return false, nil
}
