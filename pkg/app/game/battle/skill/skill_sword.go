package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type sword struct {
	ID   string
	Type int
	Arg  Argument

	count  int
	drawer skilldraw.DrawSword
}

func newSword(objID string, swordType int, arg Argument) *sword {
	res := &sword{
		ID:   objID,
		Type: swordType,
		Arg:  arg,
	}

	res.drawer.Init() // TODO error

	return res
}

func (p *sword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(p.Type, view, p.count)
}

func (p *sword) Process() (bool, error) {
	p.count++

	if p.count == 1*resources.SkillSwordDelay {
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
		case resources.SkillTypeSword:
			// No more damage area
		case resources.SkillTypeWideSword:
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
		case resources.SkillTypeLongSword:
			targetPos.X = userPos.X + 2
			if objID := field.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				localanim.DamageManager().New(dm)
			}
		}
	}

	if p.count > resources.SkillSwordEndCount {
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
