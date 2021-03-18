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
	// SkillCannon ...
	SkillCannon int = iota

	skillMax
)

const (
	typeNormalCannon int = iota
	typeHighCannon
	typeMegaCannon
	// TODO: typeWideSword, ...
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 5
)

var (
	imgCannonAtk  [3][]int32
	imgCannonBody [3][]int32
)

type cannon struct {
	Type       int
	OwnerID    string
	Power      int
	TargetType int

	count int
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
}

// Get ...
func Get(skillID int, ownerID string, targetType int) anim.Anim {
	switch skillID {
	case SkillCannon:
		return &cannon{OwnerID: ownerID, Type: typeNormalCannon, Power: 40, TargetType: targetType}
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
