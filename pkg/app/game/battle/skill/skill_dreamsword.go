package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type dreamSword struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func newDreamSword(objID string, arg Argument) *dreamSword {
	return &dreamSword{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
	}
}

func (p *dreamSword) Draw() {
	pos := objanim.GetObjPos(p.OwnerID)
	view := battlecommon.ViewPos(pos)

	n := (p.count - 5) / delaySword
	if n >= 0 && n < len(imgDreamSword) {
		dxlib.DrawRotaGraph(view.X+100, view.Y, 1, 0, imgDreamSword[n], true)
	}
}

func (p *dreamSword) Process() (bool, error) {
	p.count++

	if p.count == 1*delaySword {
		sound.On(sound.SEDreamSword)

		for x := 1; x <= 2; x++ {
			for y := -1; y <= 1; y++ {
				pos := objanim.GetObjPos(p.OwnerID)
				dm := damage.Damage{
					Power:         int(p.Power),
					TTL:           1,
					TargetType:    p.TargetType,
					HitEffectType: effect.TypeNone,
					BigDamage:     true,
					Pos:           common.Point{X: pos.X + x, Y: pos.Y + y},
				}
				damage.New(dm)
			}
		}
	}

	if p.count > len(imgDreamSword)*delaySword {
		return true, nil
	}
	return false, nil
}

func (p *dreamSword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
