package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
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
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		userPos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)

		targetPos := common.Point{X: userPos.X + 1, Y: userPos.Y}
		if objID := field.GetPanelInfo(targetPos).ObjectID; objID != "" {
			dm.TargetObjID = objID
			localanim.DamageManager().New(dm)
		}

		switch p.Type {
		case TypeSword:
			// No more damage area
		case TypeWideSword:
			targetPos.Y = userPos.Y - 1
			if objID := field.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				localanim.DamageManager().New(dm)
			}
			targetPos.Y = userPos.Y + 1
			if objID := field.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				localanim.DamageManager().New(dm)
			}
		case TypeLongSword:
			targetPos.X = userPos.X + 2
			if objID := field.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				localanim.DamageManager().New(dm)
			}
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
