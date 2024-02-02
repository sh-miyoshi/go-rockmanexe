package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type wideShot struct {
	ID   string
	Arg  Argument
	Core *processor.WideShot
}

func newWideShot(arg Argument, core skillcore.SkillCore) *wideShot {
	return &wideShot{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.WideShot),
	}
}

func (p *wideShot) Draw() {
	// nothing to do at router
}

func (p *wideShot) Process() (bool, error) {
	return p.Core.Process()
}

func (p *wideShot) GetParam() anim.Param {
	pm := p.Core.GetParam()
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeWideShot,
		ActCount:      pm.State*1000 + p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *wideShot) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *wideShot) GetEndCount() int {
	return 0
}
