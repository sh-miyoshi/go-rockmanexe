package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type recover struct {
	ID  string
	Arg Argument

	count  int
	drawer skilldraw.DrawRecover
}

func newRecover(objID string, arg Argument) *recover {
	res := &recover{
		ID:  objID,
		Arg: arg,
	}
	res.drawer.Init() // TODO: error

	return res
}

func (p *recover) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.count)
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(resources.SERecover)
		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         -int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			Element:       damage.ElementNone,
			TargetObjID:   p.Arg.OwnerID,
		})
	}

	p.count++

	if p.count > resources.SkillRecoverEndCount {
		return true, nil
	}
	return false, nil
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *recover) StopByOwner() {
	// Nothing to do
}
