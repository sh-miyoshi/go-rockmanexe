package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type shockWave struct {
	ID   string
	Arg  Argument
	Core *processor.ShockWave
}

func newShockWave(arg Argument, core skillcore.SkillCore) *shockWave {
	return &shockWave{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.ShockWave),
	}
}

func (p *shockWave) Draw() {
	// nothing to do at router
}

func (p *shockWave) Process() (bool, error) {
	return p.Core.Process()
}

func (p *shockWave) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeShockWave,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *shockWave) StopByOwner() {
}

func (p *shockWave) GetEndCount() int {
	return p.Core.GetEndCount()
}
