package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
)

type quickGauge struct {
	ID  string
	Arg Argument
}

func newQuickGauge(objID string, arg Argument) *quickGauge {
	return &quickGauge{
		ID:  objID,
		Arg: arg,
	}
}

func (p *quickGauge) Draw() {
}

func (p *quickGauge) Process() (bool, error) {
	battlecommon.CustomGaugeSpeed = 6
	return true, nil
}

func (p *quickGauge) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *quickGauge) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
