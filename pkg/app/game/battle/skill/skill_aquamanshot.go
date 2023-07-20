package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type aquamanShot struct {
	ID  string
	Arg Argument

	pos    common.Point
	ofs    common.Point
	target common.Point
	count  int
}

func newAquamanShot(objID string, arg Argument) *aquamanShot {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	return &aquamanShot{
		ID:     objID,
		Arg:    arg,
		pos:    common.Point{X: view.X - 40, Y: view.Y + 10},
		target: common.Point{X: pos.X - 2, Y: pos.Y},
	}
}

func (p *aquamanShot) Draw() {
	dxlib.DrawRotaGraph(p.pos.X+p.ofs.X, p.pos.Y+p.ofs.Y, 1, 0, imgAquamanShot[0], true)
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
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           20,
			TargetType:    p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeWater,
		})
		target := common.Point{X: p.target.X - 1, Y: p.target.Y}
		localanim.AnimNew(effect.Get(resources.EffectTypeWaterBomb, target, 0))
		localanim.DamageManager().New(damage.Damage{
			Pos:           target,
			Power:         int(p.Arg.Power),
			TTL:           20,
			TargetType:    p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeWater,
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
