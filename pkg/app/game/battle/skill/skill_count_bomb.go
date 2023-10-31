package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

type countBomb struct {
	ID  string
	Arg Argument

	count int
}

func newCountBomb(objID string, arg Argument) *countBomb {
	return &countBomb{
		ID:  objID,
		Arg: arg,
	}
}

func (p *countBomb) Draw() {
	draw.String(100, 100, 0xff0000, "TODO: 落ちてくるアニメーション")
}

func (p *countBomb) Process() (bool, error) {
	if p.count == 0 {
		field.SetBlackoutCount(150)
		SetChipNameDraw("カウントボム", true)
	}

	p.count++
	return false, nil
}

func (p *countBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *countBomb) StopByOwner() {
}
