package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type flamePillarManager struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.FlamePillarManager

	drawer skilldraw.DrawFlamePillerManager
}

func newFlamePillar(objID string, arg skillcore.Argument, core skillcore.SkillCore) *flamePillarManager {
	return &flamePillarManager{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.FlamePillarManager),
	}
}

func (p *flamePillarManager) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	// 自分の攻撃だったらBodyを表示
	showBody := p.Arg.TargetType == damage.TargetEnemy
	p.drawer.Draw(view, p.Core.GetCount(), showBody, p.Core.GetPillars(), p.Core.GetDelay())
}

func (p *flamePillarManager) Process() (bool, error) {
	return p.Core.Process()
}

func (p *flamePillarManager) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *flamePillarManager) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
