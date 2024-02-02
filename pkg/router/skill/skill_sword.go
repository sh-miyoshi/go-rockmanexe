package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type sword struct {
	ID      string
	SkillID int
	Arg     Argument
	Core    *processor.Sword
}

func newSword(skillID int, arg Argument, core skillcore.SkillCore) *sword {
	return &sword{
		ID:      arg.AnimObjID,
		SkillID: skillID,
		Arg:     arg,
		Core:    core.(*processor.Sword),
	}
}

func (p *sword) Draw() {
	// nothing to do at router
}

func (p *sword) Process() (bool, error) {
	return p.Core.Process()
}

func (p *sword) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	switch p.SkillID {
	case resources.SkillSword:
		info.AnimType = routeranim.TypeSword
	case resources.SkillWideSword:
		info.AnimType = routeranim.TypeWideSword
	case resources.SkillLongSword:
		info.AnimType = routeranim.TypeLongSword
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *sword) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *sword) GetEndCount() int {
	return p.Core.GetEndCount()
}
