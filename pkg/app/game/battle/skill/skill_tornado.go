package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	hitNum      = 8
	atkInterval = 4
)

type tornado struct {
	ID  string
	Arg skillcore.Argument

	count     int
	atkCount  int
	drawer    skilldraw.DrawTurnado
	objPos    point.Point
	targetPos point.Point
}

func newTornado(objID string, arg skillcore.Argument) *tornado {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	return &tornado{
		ID:        objID,
		Arg:       arg,
		objPos:    pos,
		targetPos: point.Point{X: pos.X + 2, Y: pos.Y},
	}
}

func (p *tornado) Draw() {
	view := battlecommon.ViewPos(p.objPos)
	target := battlecommon.ViewPos(p.targetPos)
	p.drawer.Draw(view, target, p.count)
}

func (p *tornado) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(resources.SETornado)
	}

	if p.count%atkInterval == 0 {
		lastAtk := p.atkCount == hitNum-1
		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			BigDamage:     lastAtk,
			Element:       damage.ElementNone,
			Pos:           p.targetPos,
			TTL:           atkInterval,
			ShowHitArea:   false,
		})

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
