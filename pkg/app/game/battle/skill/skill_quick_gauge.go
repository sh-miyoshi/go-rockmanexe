package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

type quickGauge struct {
	ID  string
	Arg Argument

	count int
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
	if p.count == 0 {
		field.SetBlackoutCount(90)
		SetChipNameDraw("クイックゲージ", true)
		battlecommon.CustomGaugeSpeed = 6
	}

	p.count++
	return p.count >= 90, nil
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
