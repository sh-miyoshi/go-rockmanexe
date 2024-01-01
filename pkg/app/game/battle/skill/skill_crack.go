package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	crackType1 int = iota
	crackType2
	crackType3
)

type crack struct {
	ID  string
	Arg skillcore.Argument

	count     int
	attackPos []point.Point
}

func newCrack(objID string, crackType int, arg skillcore.Argument) *crack {
	res := crack{
		ID:  objID,
		Arg: arg,
	}

	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)

	switch crackType {
	case crackType1:
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
	case crackType2:
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 2, Y: pos.Y})
	case crackType3:
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 1, Y: pos.Y - 1})
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
		res.attackPos = append(res.attackPos, point.Point{X: pos.X + 1, Y: pos.Y + 1})
	}

	return &res
}

func (p *crack) Draw() {
}

func (p *crack) Process() (bool, error) {
	p.count++

	if p.count > 5 {
		sound.On(resources.SEPanelBreak)
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
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *crack) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
