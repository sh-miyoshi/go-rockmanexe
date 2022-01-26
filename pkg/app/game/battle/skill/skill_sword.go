package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	TypeSword int = iota
	TypeWideSword
	TypeLongSword

	TypeSwordMax
)

type sword struct {
	ID         string
	Type       int
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func newSword(objID string, swordType int, arg Argument) *sword {
	return &sword{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Type:       swordType,
		Power:      arg.Power,
		TargetType: arg.TargetType,
	}
}

func (p *sword) Draw() {
	pos := objanim.GetObjPos(p.OwnerID)
	view := battlecommon.ViewPos(pos)

	n := (p.count - 5) / delaySword
	if n >= 0 && n < len(imgSword[p.Type]) {
		dxlib.DrawRotaGraph(view.X+100, view.Y, 1, 0, imgSword[p.Type][n], true)
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
			BigDamage:     true,
		}

		pos := objanim.GetObjPos(p.OwnerID)

		dm.Pos.X = pos.X + 1
		dm.Pos.Y = pos.Y
		damage.New(dm)

		switch p.Type {
		case TypeSword:
			// No more damage area
		case TypeWideSword:
			dm.Pos.Y = pos.Y - 1
			damage.New(dm)
			dm.Pos.Y = pos.Y + 1
			damage.New(dm)
		case TypeLongSword:
			dm.Pos.X = pos.X + 2
			damage.New(dm)
		}
	}

	if p.count > len(imgSword[p.Type])*delaySword {
		return true, nil
	}
	return false, nil
}

func (p *sword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
