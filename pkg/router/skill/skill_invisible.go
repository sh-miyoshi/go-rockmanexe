package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type invisible struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func newInvisible(arg Argument, core skillcore.SkillCore) *invisible {
	return &invisible{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *invisible) Draw() {
	// nothing to do at router
}

func (p *invisible) Update() (bool, error) {
	return p.Core.Update()
}

func (p *invisible) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeInvisible,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *invisible) StopByOwner() {
}
