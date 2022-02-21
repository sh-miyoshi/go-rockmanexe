package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

type invisible struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func newInvisible(objID string, arg Argument) *invisible {
	return &invisible{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
	}
}

func (p *invisible) Draw() {
}

func (p *invisible) Process() (bool, error) {
	p.count++

	showTm := 60
	if p.count == 1 {
		field.SetBlackoutCount(showTm)
		objanim.MakeInvisible(p.OwnerID, 6*60)
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
