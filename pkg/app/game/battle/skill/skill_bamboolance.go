package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type bambooLance struct {
	ID      string
	Arg     skillcore.Argument
	Core    skillcore.SkillCore
	drawer  skilldraw.DrawBamboolance
	animMgr *manager.Manager
}

func newBambooLance(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *bambooLance {
	res := &bambooLance{
		ID:      objID,
		Arg:     arg,
		Core:    core,
		animMgr: animMgr,
	}
	res.drawer.Init()

	return res
}

func (p *bambooLance) Draw() {
	p.drawer.Draw(p.Core.GetCount(), true)
}

func (p *bambooLance) Update() (bool, error) {
	return p.Core.Update()
}

func (p *bambooLance) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *bambooLance) StopByOwner() {
	// Nothing to do after throwing
}
