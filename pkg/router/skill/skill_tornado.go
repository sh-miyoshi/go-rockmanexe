package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type tornado struct {
	ID   string
	Arg  Argument
	Core *processor.Tornado
}

func newTornado(arg Argument, core skillcore.SkillCore) *tornado {
	return &tornado{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.Tornado),
	}
}

func (p *tornado) Draw() {
	// nothing to do at router
}

func (p *tornado) Process() (bool, error) {
	return p.Core.Process()
}

func (p *tornado) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeTornado,
	}

	obj, _ := p.Core.GetPos()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       obj,
		ExtraInfo: info.Marshal(),
	}
}

func (p *tornado) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}

func (p *tornado) GetEndCount() int {
	return p.Core.GetEndCount()
}
