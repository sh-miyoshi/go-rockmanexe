package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillForteShootingBuster struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteShootingBuster

	drawer  skilldraw.DrawForteShootingBuster
	animMgr *manager.Manager
}

func newForteShootingBuster(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *skillForteShootingBuster {
	return &skillForteShootingBuster{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ForteShootingBuster),
		animMgr: animMgr,
	}
}

func (p *skillForteShootingBuster) Draw() {
	p.drawer.Draw(p.Core.GetPos(), p.Core.GetCount(), p.Core.GetInitWait())
}

func (p *skillForteShootingBuster) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillForteShootingBuster) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *skillForteShootingBuster) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
