package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	miniBombEndCount   = 60
	delayMiniBombThrow = 4
)

type miniBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    common.Point
	target common.Point
}

func newMiniBomb(objID string, arg Argument) *miniBomb {
	pos := objanim.GetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:     objID,
		Arg:    arg,
		pos:    pos,
		target: common.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	imgNo := (p.count / delayMiniBombThrow) % len(imgBombThrow)
	view := battlecommon.ViewPos(p.pos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := battlecommon.PanelSize.X * (p.target.X - p.pos.X)
	ofsx := size * p.count / miniBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.target.Y != p.pos.Y {
		size = battlecommon.PanelSize.Y * (p.target.Y - p.pos.Y)
		dy := size * p.count / miniBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgBombThrow[imgNo], true)
}

func (p *miniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEBombThrow)
	}

	if p.count == miniBombEndCount {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEExplode)
		anim.New(effect.Get(effect.TypeExplode, p.target, 0))
		damage.New(damage.Damage{
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		})
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *miniBomb) StopByOwner() {
	if p.count < 5 {
		anim.Delete(p.ID)
	}
}
