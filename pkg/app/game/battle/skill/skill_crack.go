package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	crackType1 int = iota
	crackType2
	crackType3
)

type crack struct {
	ID  string
	Arg Argument

	count     int
	attackPos []common.Point
}

func newCrack(objID string, crackType int, arg Argument) *crack {
	res := crack{
		ID:  objID,
		Arg: arg,
	}

	pos := objanim.GetObjPos(arg.OwnerID)

	switch crackType {
	case crackType1:
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 1, Y: pos.Y})
	case crackType2:
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 1, Y: pos.Y})
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 2, Y: pos.Y})
	case crackType3:
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 1, Y: pos.Y - 1})
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 1, Y: pos.Y})
		res.attackPos = append(res.attackPos, common.Point{X: pos.X + 1, Y: pos.Y + 1})
	}

	return &res
}

func (p *crack) Draw() {
}

func (p *crack) Process() (bool, error) {
	p.count++

	if p.count > 5 {
		sound.On(sound.SEPanelBreak)
		for _, pos := range p.attackPos {
			field.PanelBreak(pos)
		}

		return true, nil
	}

	return false, nil
}

func (p *crack) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *crack) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}
