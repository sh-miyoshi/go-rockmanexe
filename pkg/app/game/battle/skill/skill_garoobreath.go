package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
)

type garooBreath struct {
	ID  string
	Arg Argument

	count    int
	pos      common.Point
	next     common.Point
	prev     common.Point
	moveVecX int
	damageID string
	drawer   skilldraw.DrawGarooBreath
}

func newGarooBreath(objID string, arg Argument) *garooBreath {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	vx := 1
	if arg.TargetType == damage.TargetPlayer {
		vx = -1
	}
	first := common.Point{X: pos.X + vx, Y: pos.Y}

	return &garooBreath{
		ID:       objID,
		Arg:      arg,
		pos:      first,
		prev:     pos,
		next:     first,
		moveVecX: vx,
	}
}

func (p *garooBreath) Draw() {
	p.drawer.Draw(p.prev, p.pos, p.next, p.count)
}

func (p *garooBreath) Process() (bool, error) {
	if p.count%resources.SkillGarooBreathNextStepCount/2 == 0 {
		// Finish if hit
		if p.damageID != "" && !localanim.DamageManager().Exists(p.damageID) {
			return true, nil
		}
	}

	if p.count%resources.SkillGarooBreathNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		p.damageID = localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           resources.SkillGarooBreathNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeHeatHit,
			ShowHitArea:   false,
			BigDamage:     true,
			Element:       damage.ElementFire,
		})

		// Set next pos
		p.next.X += p.moveVecX
	}

	p.count++

	if p.pos.X < 0 || p.pos.X >= battlecommon.FieldNum.X {
		return true, nil
	}
	return false, nil
}

func (p *garooBreath) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *garooBreath) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
