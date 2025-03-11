package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type cannon struct {
	ID   string
	Arg  Argument
	Core *processor.Cannon
}

func newCannon(arg Argument, core skillcore.SkillCore) *cannon {
	return &cannon{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.Cannon),
	}
}

func (p *cannon) Draw() {
	// nothing to do at router
}

func (p *cannon) Update() (bool, error) {
	return p.Core.Update()
}

func (p *cannon) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID: p.ID,

		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *cannon) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}
