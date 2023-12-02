package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type miniBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    point.Point
	target point.Point
	drawer skilldraw.DrawMiniBomb
}

func newMiniBomb(objID string, arg Argument) *miniBomb {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:     objID,
		Arg:    arg,
		pos:    pos,
		target: point.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	p.drawer.Draw(p.pos, p.target, p.count)
}

func (p *miniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(resources.SEBombThrow)
	}

	if p.count == resources.SkillMiniBombEndCount {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(resources.SEExplode)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, p.target, 0))
		if objID := field.GetPanelInfo(p.target).ObjectID; objID != "" {
			localanim.DamageManager().New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementNone,
				TargetObjID:   objID,
			})
		}
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *miniBomb) StopByOwner() {
	if p.count < 5 {
		localanim.AnimDelete(p.ID)
	}
}
