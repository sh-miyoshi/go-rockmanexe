package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type cannon struct {
	ID   string
	Type int
	Arg  Argument

	count  int
	drawer skilldraw.DrawCannon
}

func newCannon(objID string, cannonType int, arg Argument) *cannon {
	return &cannon{
		ID:   objID,
		Type: cannonType,
		Arg:  arg,
	}
}

func (p *cannon) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(p.Type, view, p.count)
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		sound.On(resources.SECannon)
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		dm := damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeCannonHit,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				if objID := field.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					localanim.DamageManager().New(dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				if objID := field.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					localanim.DamageManager().New(dm)
					break
				}
			}
		}
	}

	if p.count > resources.SkillCannonEndCount {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *cannon) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
