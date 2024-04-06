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

type heatShot struct {
	ID      string
	SkillID int
	Arg     Argument
	Core    *processor.HeatShot
}

func newHeatShot(skillID int, arg Argument, core skillcore.SkillCore) *heatShot {
	return &heatShot{
		ID:      arg.AnimObjID,
		SkillID: skillID,
		Arg:     arg,
		Core:    core.(*processor.HeatShot),
	}
}

func (p *heatShot) Draw() {
	// nothing to do at router
}

func (p *heatShot) Process() (bool, error) {
	res, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		p.Arg.Manager.QueuePush(gameinfo.QueueTypeEffect, &gameinfo.Effect{
			ID:            uuid.New().String(),
			OwnerClientID: p.Arg.OwnerClientID,
			Pos:           hit,
			Type:          resources.EffectTypeHeatHit,
			RandRange:     0,
		})
	}

	return res, nil
}

func (p *heatShot) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	switch p.SkillID {
	case resources.SkillHeatShot:
		info.AnimType = routeranim.TypeHeatShot
	case resources.SkillHeatV:
		info.AnimType = routeranim.TypeHeatV
	case resources.SkillHeatSide:
		info.AnimType = routeranim.TypeHeatSide
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *heatShot) StopByOwner() {
	if p.Core.GetCount() < p.Core.GetDelay() {
		p.Arg.Manager.AnimDelete(p.ID)
	}
}
