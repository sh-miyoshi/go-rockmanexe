package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type waterBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    common.Point
	target common.Point
	drawer skilldraw.DrawWaterBomb
}

func newWaterBomb(objID string, arg Argument) *waterBomb {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	t := common.Point{X: pos.X + 3, Y: pos.Y}
	objType := objanim.ObjTypePlayer
	if arg.TargetType == damage.TargetEnemy {
		objType = objanim.ObjTypeEnemy
	}

	objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objType})
	if len(objs) > 0 {
		t = objs[0].Pos
	}

	return &waterBomb{
		ID:     objID,
		Arg:    arg,
		target: t,
		pos:    pos,
	}
}

func (p *waterBomb) Draw() {
	p.drawer.Draw(p.pos, p.target, p.count)
}

func (p *waterBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(resources.SEBombThrow)
	}

	if p.count == resources.SkillWaterBombEndCount {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(resources.SEWaterLanding)
		localanim.AnimNew(effect.Get(resources.EffectTypeWaterBomb, p.target, 0))
		if objID := field.GetPanelInfo(p.target).ObjectID; objID != "" {
			localanim.DamageManager().New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementWater,
				TargetObjID:   objID,
			})
		}
		field.PanelCrack(p.target)
		return true, nil
	}
	return false, nil
}

func (p *waterBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *waterBomb) StopByOwner() {
	if p.count < 5 {
		localanim.AnimDelete(p.ID)
	}
}
