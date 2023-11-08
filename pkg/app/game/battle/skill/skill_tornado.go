package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
)

const (
	hitNum      = 8
	atkInterval = 4
)

type tornado struct {
	ID  string
	Arg Argument

	count    int
	atkCount int
	drawer   skilldraw.DrawTurnado
}

func newTornado(objID string, arg Argument) *tornado {
	return &tornado{
		ID:  objID,
		Arg: arg,
	}
}

func (p *tornado) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	targetPos := common.Point{X: pos.X + 2, Y: pos.Y}
	view := battlecommon.ViewPos(pos)
	target := battlecommon.ViewPos(targetPos)
	p.drawer.Draw(view, target, p.count)
}

func (p *tornado) Process() (bool, error) {
	p.count++

	if p.count%atkInterval == 0 {
		// TODO: ダメージ処理

		p.atkCount++
		return p.atkCount >= hitNum, nil
	}

	return false, nil
}

func (p *tornado) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *tornado) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
