package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type miniBomb struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore

	pos    point.Point
	target point.Point
	drawer skilldraw.DrawMiniBomb
}

func newMiniBomb(objID string, arg skillcore.Argument, core skillcore.SkillCore) *miniBomb {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:     objID,
		Arg:    arg,
		Core:   core,
		pos:    pos,
		target: point.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	p.drawer.Draw(p.pos, p.target, p.Core.GetCount())
}

func (p *miniBomb) Process() (bool, error) {
	end, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	if end {
		sound.On(resources.SEExplode)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, p.target, 0))
	}
	return end, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *miniBomb) StopByOwner() {
	if p.Core.GetCount() < 5 {
		localanim.AnimDelete(p.ID)
	}
}
