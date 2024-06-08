package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type shrimpyAtk struct {
	ID  string
	Arg skillcore.Argument

	drawer skilldraw.DrawShrimpyAtk
}

func newShrimpyAtk(objID string, arg skillcore.Argument) *shrimpyAtk {
	return &shrimpyAtk{
		ID:  objID,
		Arg: arg,
	}
}

func (p *shrimpyAtk) Draw() {
	p.drawer.Draw()
}

func (p *shrimpyAtk) Process() (bool, error) {
	// TODO
	return false, nil
}

func (p *shrimpyAtk) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *shrimpyAtk) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
