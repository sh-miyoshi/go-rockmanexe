package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type ComeOnSnake struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.ComeOnSnake
	animMgr *manager.Manager
}

func newComeOnSnake(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *ComeOnSnake {
	return &ComeOnSnake{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ComeOnSnake),
		animMgr: animMgr,
	}
}

func (p *ComeOnSnake) Draw() {
	// TODO: implement draw method
}

func (p *ComeOnSnake) Update() (bool, error) {
	return p.Core.Update()
}

func (p *ComeOnSnake) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *ComeOnSnake) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
