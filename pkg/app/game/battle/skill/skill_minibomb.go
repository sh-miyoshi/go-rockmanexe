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
	miniBombEndCount = 60
)

type miniBomb struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count  int
	pos    common.Point
	target common.Point
}

func newMiniBomb(objID string, arg Argument) *miniBomb {
	pos := objanim.GetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		pos:        pos,
		target:     common.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	imgNo := (p.count / delayBombThrow) % len(imgBombThrow)
	view := battlecommon.ViewPos(p.pos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := field.PanelSize.X * (p.target.X - p.pos.X)
	ofsx := size * p.count / miniBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.target.Y != p.pos.Y {
		size = field.PanelSize.Y * (p.target.Y - p.pos.Y)
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
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEExplode)
		anim.New(effect.Get(effect.TypeExplode, p.target, 0))
		damage.New(damage.Damage{
			Pos:           p.target,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
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
