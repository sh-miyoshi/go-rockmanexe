package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

type invisible struct {
	ID  string
	Arg Argument

	count int
}

func newInvisible(objID string, arg Argument) *invisible {
	return &invisible{
		ID:  objID,
		Arg: arg,
	}
}

func (p *invisible) Draw() {
}

func (p *invisible) Process() (bool, error) {
	p.count++

	showTm := 60
	if p.count == 1 {
		field.SetBlackoutCount(showTm)
		objanim.MakeInvisible(p.Arg.OwnerID, 6*60)
		setChipNameDraw("インビジブル")
	}

	return p.count > showTm, nil
}

func (p *invisible) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *invisible) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}
