package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type recover struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func newRecover(arg Argument, core skillcore.SkillCore) *recover {
	return &recover{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *recover) Draw() {
	// nothing to do at router
}

func (p *recover) Update() (bool, error) {
	return p.Core.Update()
}

func (p *recover) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeRecover,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *recover) StopByOwner() {
}
