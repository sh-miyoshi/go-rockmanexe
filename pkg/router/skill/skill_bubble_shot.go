package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type bubbleShot struct {
	ID   string
	Arg  Argument
	Core *processor.BubbleShot
}

func newBubbleShot(arg Argument, core skillcore.SkillCore) *bubbleShot {
	return &bubbleShot{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.BubbleShot),
	}
}

func (p *bubbleShot) Draw() {
	// nothing to do at router
}

func (p *bubbleShot) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		p.Arg.Manager.QueuePush(gameinfo.QueueTypeEffect, &gameinfo.Effect{
			ID:            uuid.New().String(),
			OwnerClientID: p.Arg.OwnerClientID,
			Pos:           hit,
			Type:          resources.EffectTypeWaterBomb,
			RandRange:     0,
		})
	}
	return res, nil
}

func (p *bubbleShot) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeBubbleShot,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *bubbleShot) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}
