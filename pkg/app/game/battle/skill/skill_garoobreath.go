package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	garooBreathNextStepCount = 10
	delayGarooBreath         = 4
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
}

func newGarooBreath(objID string, arg Argument) *garooBreath {
	pos := objanim.GetObjPos(arg.OwnerID)
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
	view := battlecommon.ViewPos(p.pos)
	n := (p.count / delayGarooBreath) % len(imgGarooBreath)

	cnt := p.count % garooBreathNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(p.next.X, p.pos.X, p.prev.X, cnt, garooBreathNextStepCount, battlecommon.PanelSize.X)
	ofsy := -15
	xflip := int32(dxlib.TRUE)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgGarooBreath[n], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
}

func (p *garooBreath) Process() (bool, error) {
	if p.count%garooBreathNextStepCount/2 == 0 {
		// Finish if hit
		if p.damageID != "" && !damage.Exists(p.damageID) {
			return true, nil
		}
	}

	if p.count%garooBreathNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		p.damageID = damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           garooBreathNextStepCount + 1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeHeatHit,
			ShowHitArea:   false,
			BigDamage:     true,
			DamageType:    damage.TypeFire,
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
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *garooBreath) StopByOwner() {
	anim.Delete(p.ID)
}
