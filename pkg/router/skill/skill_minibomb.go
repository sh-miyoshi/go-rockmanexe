package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type miniBomb struct {
	ID   string
	Arg  Argument
	Core skillcore.SkillCore

	pos    point.Point
	target point.Point
}

func newMiniBomb(arg Argument, core skillcore.SkillCore) *miniBomb {
	pos := routeranim.ObjAnimGetObjPos(arg.OwnerClientID, arg.OwnerObjectID)
	return &miniBomb{
		ID:     arg.AnimObjID,
		Arg:    arg,
		Core:   core,
		pos:    pos,
		target: point.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	// nothing to do at router
}

func (p *miniBomb) Process() (bool, error) {
	return p.Core.Process()
}

func (p *miniBomb) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeMiniBomb,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.pos,
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *miniBomb) StopByOwner() {
	if p.Core.GetCount() < 5 {
		routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
	}
}

func (p *miniBomb) GetEndCount() int {
	return p.Core.GetEndCount()
}
