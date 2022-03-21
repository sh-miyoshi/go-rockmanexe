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
	waterBombEndCount   = 60
	delayWaterBombThrow = 4
)

type waterBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    common.Point
	target common.Point
}

func newWaterBomb(objID string, arg Argument) *waterBomb {
	pos := objanim.GetObjPos(arg.OwnerID)
	t := common.Point{X: pos.X + 3, Y: pos.Y}
	objType := objanim.ObjTypePlayer
	if arg.TargetType == damage.TargetEnemy {
		objType = objanim.ObjTypeEnemy
	}

	objs := objanim.GetObjs(objanim.Filter{ObjType: objType})
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
	imgNo := (p.count / delayWaterBombThrow) % len(imgBombThrow)
	view := battlecommon.ViewPos(p.pos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := field.PanelSize.X * (p.target.X - p.pos.X)
	ofsx := size * p.count / waterBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.target.Y != p.pos.Y {
		size = field.PanelSize.Y * (p.target.Y - p.pos.Y)
		dy := size * p.count / waterBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgBombThrow[imgNo], true)
}

func (p *waterBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEBombThrow)
	}

	if p.count == waterBombEndCount {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEWaterLanding)
		anim.New(effect.Get(effect.TypeWaterBomb, p.target, 0))
		damage.New(damage.Damage{
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})
		field.PanelCrack(p.target)
		return true, nil
	}
	return false, nil
}

func (p *waterBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *waterBomb) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}
