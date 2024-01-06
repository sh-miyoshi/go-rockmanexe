package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type invisible struct {
	ID  string
	Arg skillcore.Argument

	count int
}

func newInvisible(objID string, arg skillcore.Argument) *invisible {
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
		localanim.ObjAnimMakeInvisible(p.Arg.OwnerID, 6*60)
		SetChipNameDraw("インビジブル", true)
	}

	return p.count > showTm, nil
}

func (p *invisible) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *invisible) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
