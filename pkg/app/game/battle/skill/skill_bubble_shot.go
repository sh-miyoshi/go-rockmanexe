package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type bubbleShot struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.BubbleShot

	drawer skilldraw.DrawBubbleShot
}

func newBubbleShot(objID string, arg skillcore.Argument, core skillcore.SkillCore) *bubbleShot {
	return &bubbleShot{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.BubbleShot),
	}
}

func (p *bubbleShot) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), true)
}

func (p *bubbleShot) Process() (bool, error) {
	res, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		localanim.AnimNew(effect.Get(resources.EffectTypeWaterBomb, hit, 0))
	}
	return res, nil
}

func (p *bubbleShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *bubbleShot) StopByOwner() {
	if p.Core.GetCount() < p.Core.GetDelay() {
		localanim.AnimDelete(p.ID)
	}
}