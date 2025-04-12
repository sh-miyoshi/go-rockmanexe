package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
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

	drawer   skilldraw.DrawBubbleShot
	animMgr  *manager.Manager
	showAnim bool
}

func newBubbleShot(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager, showAnim bool) *bubbleShot {
	return &bubbleShot{
		ID:       objID,
		Arg:      arg,
		Core:     core.(*processor.BubbleShot),
		animMgr:  animMgr,
		showAnim: showAnim,
	}
}

func (p *bubbleShot) Draw() {
	if !p.showAnim {
		return
	}

	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), true)
}

func (p *bubbleShot) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		p.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeWaterBomb, hit, 0))
	}
	return res, nil
}

func (p *bubbleShot) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *bubbleShot) StopByOwner() {
	if p.Core.GetCount() < p.Core.GetDelay() {
		p.animMgr.AnimDelete(p.ID)
	}
}
