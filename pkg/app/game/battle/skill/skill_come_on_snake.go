package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type ComeOnSnake struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.ComeOnSnake
	animMgr *manager.Manager
	drawer  skilldraw.DrawSnake
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
	for _, snake := range p.Core.GetSnakes() {
		p.drawer.Draw(snake.ViewPos, snake.Count)
	}
}

func (p *ComeOnSnake) Update() (bool, error) {
	end, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		return true, nil
	}
	return false, nil
}

func (p *ComeOnSnake) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *ComeOnSnake) StopByOwner() {
	// Nothing to do
}
