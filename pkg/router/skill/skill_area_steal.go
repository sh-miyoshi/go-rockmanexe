package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type areaSteal struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func newAreaSteal(arg Argument, core skillcore.SkillCore) *areaSteal {
	return &areaSteal{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *areaSteal) Draw() {
	// nothing to do at router
}

func (p *areaSteal) Process() (bool, error) {
	return p.Core.Process()
}

func (p *areaSteal) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeAreaSteal,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *areaSteal) StopByOwner() {
}
