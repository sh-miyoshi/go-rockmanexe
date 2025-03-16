package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type crack struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func newCrack(arg Argument, core skillcore.SkillCore) *crack {
	return &crack{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *crack) Draw() {
	// nothing to do at router
}

func (p *crack) Update() (bool, error) {
	return p.Core.Update()
}

func (p *crack) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeCrack,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *crack) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}
