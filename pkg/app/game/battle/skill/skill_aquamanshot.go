package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type aquamanShot struct {
	ID  string
	Arg Argument

	pos    common.Point
	ofs    common.Point
	target common.Point
	count  int
	drawer skilldraw.DrawAquamanShot
}

func newAquamanShot(objID string, arg Argument) *aquamanShot {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	res := &aquamanShot{
		ID:     objID,
		Arg:    arg,
		pos:    common.Point{X: view.X - 40, Y: view.Y + 10},
		target: common.Point{X: pos.X - 2, Y: pos.Y},
	}
	res.drawer.Init()

	return res
}

func (p *aquamanShot) Draw() {
	p.drawer.Draw(p.pos, p.ofs)
}

func (p *aquamanShot) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(resources.SEBombThrow)
	}

	const size = 120
	p.ofs.X -= 6
	p.ofs.Y = 10*p.ofs.X*p.ofs.X/(size*size) - 20*p.ofs.X/size

	if p.ofs.X < -size {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(resources.SEWaterLanding)
		localanim.AnimNew(effect.Get(resources.EffectTypeWaterBomb, p.target, 0))
		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           20,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementWater,
		})
		target := common.Point{X: p.target.X - 1, Y: p.target.Y}
		localanim.AnimNew(effect.Get(resources.EffectTypeWaterBomb, target, 0))
		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           target,
			Power:         int(p.Arg.Power),
			TTL:           20,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementWater,
		})

		return true, nil
	}
	return false, nil
}

func (p *aquamanShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *aquamanShot) StopByOwner() {
	// Nothing to do after throwing
}
