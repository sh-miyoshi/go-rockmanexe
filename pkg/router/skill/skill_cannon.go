package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	TypeNormalCannon int = iota
	TypeHighCannon
	TypeMegaCannon

	TypeCannonMax
)

type cannon struct {
	ID   string
	Type int
	Arg  Argument
	Core skillcore.SkillCore
}

func newCannon(cannonType int, arg Argument, core skillcore.SkillCore) *cannon {
	return &cannon{
		ID:   arg.AnimObjID,
		Type: cannonType,
		Arg:  arg,
		Core: core,
	}
}

func (p *cannon) Draw() {
	// nothing to do at router
}

func (p *cannon) Process() (bool, error) {
	return p.Core.Process()
}

func (p *cannon) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	switch p.Type {
	case TypeNormalCannon:
		info.AnimType = routeranim.TypeCannonNormal
	case TypeHighCannon:
		info.AnimType = routeranim.TypeCannonHigh
	case TypeMegaCannon:
		info.AnimType = routeranim.TypeCannonMega
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *cannon) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *cannon) GetEndCount() int {
	return resources.SkillCannonEndCount
}
