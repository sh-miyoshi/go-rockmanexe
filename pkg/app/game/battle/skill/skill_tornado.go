package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
)

type tornado struct {
	ID  string
	Arg Argument

	count int
}

func newTornado(objID string, arg Argument) *tornado {
	return &tornado{
		ID:  objID,
		Arg: arg,
	}
}

func (p *tornado) Draw() {
	// p.drawer.Draw()
}

func (p *tornado) Process() (bool, error) {
	p.count++

	return false, nil
}

func (p *tornado) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *tornado) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
