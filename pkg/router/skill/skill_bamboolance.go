package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type bambooLance struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore
}

func newBambooLance(arg Argument, core skillcore.SkillCore) *bambooLance {
	return &bambooLance{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core,
	}
}

func (p *bambooLance) Draw() {
	// nothing to do at router
}

func (p *bambooLance) Process() (bool, error) {
	return p.Core.Process()
}

func (p *bambooLance) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeBambooLance,
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *bambooLance) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}
