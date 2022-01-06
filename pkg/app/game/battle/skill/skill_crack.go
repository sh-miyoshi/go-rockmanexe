package skill

import (
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
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count     int
	attackPos [][]int
}

func newCrack(objID string, crackType int, arg Argument) *crack {
	res := crack{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
	}

	x, y := objanim.GetObjPos(arg.OwnerID)

	switch crackType {
	case crackType1:
		res.attackPos = append(res.attackPos, []int{x + 1, y})
	case crackType2:
		res.attackPos = append(res.attackPos, []int{x + 1, y})
		res.attackPos = append(res.attackPos, []int{x + 2, y})
	case crackType3:
		res.attackPos = append(res.attackPos, []int{x + 1, y - 1})
		res.attackPos = append(res.attackPos, []int{x + 1, y})
		res.attackPos = append(res.attackPos, []int{x + 1, y + 1})
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
			field.PanelBreak(pos[0], pos[1])
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
