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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
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
	imgNo := (p.count / delayWaterBombThrow) % len(imgBombThrow)
	view := battlecommon.ViewPos(p.pos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := battlecommon.PanelSize.X * (p.target.X - p.pos.X)
	ofsx := size * p.count / waterBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if p.target.Y != p.pos.Y {
		size = battlecommon.PanelSize.Y * (p.target.Y - p.pos.Y)
		dy := size * p.count / waterBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgBombThrow[imgNo], true)
}

func (p *waterBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(resources.SEBombThrow)
	}

	if p.count == waterBombEndCount {
		pn := field.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		sound.On(resources.SEWaterLanding)
		localanim.AnimNew(effect.Get(battlecommon.EffectTypeWaterBomb, p.target, 0))
		localanim.DamageManager().New(damage.Damage{
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: battlecommon.EffectTypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeWater,
		})
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
