package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	cirkillShotNextStepCount = 20
)

type cirkillShot struct {
	ID  string
	Arg skillcore.Argument

	count    int
	pos      point.Point
	next     point.Point
	prev     point.Point
	moveVecX int
	damageID string
	drawer   skilldraw.DrawCirkillShot
	animMgr  *manager.Manager
}

func newCirkillShot(objID string, arg skillcore.Argument, animMgr *manager.Manager) *cirkillShot {
	pos := animMgr.ObjAnimGetObjPos(arg.OwnerID)
	vx := 1
	if arg.TargetType == damage.TargetPlayer {
		vx = -1
	}
	first := point.Point{X: pos.X + vx, Y: pos.Y}

	return &cirkillShot{
		ID:       objID,
		Arg:      arg,
		pos:      first,
		prev:     pos,
		next:     first,
		moveVecX: vx,
		animMgr:  animMgr,
	}
}

func (p *cirkillShot) Draw() {
	p.drawer.Draw(p.prev, p.pos, p.next, p.count, cirkillShotNextStepCount)
}

func (p *cirkillShot) Update() (bool, error) {
	if p.count%cirkillShotNextStepCount/2 == 0 {
		// Finish if hit
		if p.damageID != "" && !p.animMgr.DamageManager().Exists(p.damageID) {
			return true, nil
		}
	}

	if p.count%cirkillShotNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		p.damageID = p.animMgr.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           cirkillShotNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeHeatHit,
			ShowHitArea:   true,
			StrengthType:  damage.StrengthHigh,
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

func (p *cirkillShot) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *cirkillShot) StopByOwner() {
}
