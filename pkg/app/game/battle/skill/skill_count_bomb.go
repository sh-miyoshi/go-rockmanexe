package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
)

const (
	endCount = 30
)

type countBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    common.Point
	drawer skilldraw.DrawCountBomb
}

func newCountBomb(objID string, arg Argument) *countBomb {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	return &countBomb{
		ID:  objID,
		Arg: arg,
		pos: common.Point{X: pos.X + 1, Y: pos.Y},
	}
}

func (p *countBomb) Draw() {
	if p.count < endCount {
		view := battlecommon.ViewPos(p.pos)
		p.drawer.Draw(view, p.count)
	}
}

func (p *countBomb) Process() (bool, error) {
	if p.count == 0 {
		field.SetBlackoutCount(endCount)
		SetChipNameDraw("カウントボム", true)
	}

	if p.count == endCount {
		obj := &object.CountBomb{}
		pm := object.ObjectParam{
			Pos:           p.pos,
			HP:            50,
			OnwerCharType: objanim.ObjTypePlayer,
			Power:         int(p.Arg.Power),
		}
		if err := obj.Init(p.ID, pm); err != nil {
			return false, fmt.Errorf("count bomb create failed: %w", err)
		}
		localanim.ObjAnimNew(obj)
		return true, nil
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
