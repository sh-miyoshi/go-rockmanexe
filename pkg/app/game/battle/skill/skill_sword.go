package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	TypeSword int = iota
	TypeWideSword
	TypeLongSword

	TypeSwordMax
)

const (
	delaySword = 3
)

type sword struct {
	ID   string
	Type int
	Arg  Argument

	count int
}

func newSword(objID string, swordType int, arg Argument) *sword {
	return &sword{
		ID:   objID,
		Type: swordType,
		Arg:  arg,
	}
}

func (p *sword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	n := (p.count - 5) / delaySword
	if n >= 0 && n < len(imgSword[p.Type]) {
		dxlib.DrawRotaGraph(view.X+100, view.Y, 1, 0, imgSword[p.Type][n], true)
	}
}

func (p *sword) Process() (bool, error) {
	p.count++

	if p.count == 1*delaySword {
		sound.On(resources.SESword)

		dm := damage.Damage{
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		}

		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)

		dm.Pos.X = pos.X + 1
		dm.Pos.Y = pos.Y
		localanim.DamageManager().New(dm)

		switch p.Type {
		case TypeSword:
			// No more damage area
		case TypeWideSword:
			dm.Pos.Y = pos.Y - 1
			localanim.DamageManager().New(dm)
			dm.Pos.Y = pos.Y + 1
			localanim.DamageManager().New(dm)
		case TypeLongSword:
			dm.Pos.X = pos.X + 2
			localanim.DamageManager().New(dm)
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
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *sword) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
