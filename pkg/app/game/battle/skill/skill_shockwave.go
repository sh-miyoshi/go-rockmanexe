package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type shockWave struct {
	ID       string
	Arg      skillcore.Argument
	ShowPick bool
	Core     (*processor.ShockWave)

	drawer     skilldraw.DrawShockWave
	pickDrawer skilldraw.DrawPick
}

func newShockWave(objID string, arg skillcore.Argument, core skillcore.SkillCore) *shockWave {
	res := &shockWave{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ShockWave),
	}

	return res
}

func (p *shockWave) Draw() {
	pm := p.Core.GetParam()
	showWave := p.Core.GetCount() > pm.InitWait
	if showWave {
		view := battlecommon.ViewPos(p.Core.GetPos())
		p.drawer.Draw(view, p.Core.GetCount(), pm.Speed, pm.Direct)
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
	if p.Core.GetCount() <= p.Core.GetParam().InitWait {
		localanim.AnimDelete(p.ID)
	}
}
