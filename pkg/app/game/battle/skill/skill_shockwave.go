package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	skilldefines "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/defines"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type shockWave struct {
	ID       string
	Arg      skillcore.Argument
	ShowPick bool
	Core     skillcore.SkillCore

	pm         skilldefines.ShockWaveParam
	pos        point.Point
	showWave   bool
	drawer     skilldraw.DrawShockWave
	pickDrawer skilldraw.DrawPick
}

func newShockWave(objID string, isPlayer bool, arg skillcore.Argument, core skillcore.SkillCore) *shockWave {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	res := &shockWave{
		ID:   objID,
		Arg:  arg,
		Core: core,
		pos:  pos,
		pm:   skilldefines.GetShockWaveParam(isPlayer),
	}

	return res
}

func (p *shockWave) Draw() {
	if p.showWave {
		view := battlecommon.ViewPos(p.pos)
		p.drawer.Draw(view, p.Core.GetCount(), p.pm.Speed, p.pm.Direct)
	}

	if p.ShowPick {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		p.pickDrawer.Draw(view, p.Core.GetCount())
	}
}

func (p *shockWave) Process() (bool, error) {
	return p.Core.Process()
}

func (p *shockWave) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *shockWave) StopByOwner() {
	if p.Core.GetCount() <= p.pm.InitWait {
		localanim.AnimDelete(p.ID)
	}
}
